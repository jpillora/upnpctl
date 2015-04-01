package main

import (
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/jpillora/upnpctl/upnp" //vendored
)

var VERSION string = "0.0.0" //set via ldflags

var helpFooter = `
	  -v, verbose logs

	Read more: https://github.com/jpillora/upnpctl
`

var help = `
	Usage: upnpctl [-v] <command>
	
	Version: ` + VERSION + `

	Commands:
	  * list: discovers all available UPnP devices
	  * add: adds a set of port mappings to a device
	  * rem: removes a set of port mappings from a device

` + helpFooter

var helpAdd = `
	Usage: upnpctl [-v] add [options] [mapping]...

	a [mapping] is an external port and optional internal
	port, which comes in the form "external[:internal}".
	for example, "3000" and "5000:6000" would be valid
	[mappings]. you may specify any number of mappings.

	Options:
	  --id, the device id. required	when more than one
	  device is found.

	  --type, port type: tcp or udp (defaults to 'tcp')

	  --timeout, port mapping timeout (defaults to permanent)

	  --desc, port mapping description; some routers
	  display this description along-side port mappings
	  (defaults to 'upnpctl v` + VERSION + `')
` + helpFooter

var helpRem = `
	Usage: upnpctl [-v] rem [options] [external]...

	a [external] is the external port identifying a port
	mapping to remove. you may specify any number of
	external ports.

	Options:
	  --id, the device id. required	when more than one
	  device is found.
` + helpFooter

type command string

var list = command("list")
var add = command("add")
var rem = command("rem")

func main() {
	v := flag.Bool("v", false, "")
	flag.Usage = func() {
		display(help)
	}
	flag.Parse()
	if *v {
		upnp.Debug = true
		upnp.EnableLog()
	}
	args := flag.Args()
	if len(args) == 0 {
		display(help)
	}

	cmd := command(args[0])
	args = args[1:]

	switch cmd {
	case list:
		listMappings()
		os.Exit(0)
	case add:
		if len(args) == 0 {
			display(helpAdd)
		}
	case rem:
		if len(args) == 0 {
			display(helpRem)
		}
	default:
		fmt.Println("no match " + cmd)
		display(help)
	}

	f := flag.NewFlagSet(string(cmd), flag.ExitOnError)
	id := f.String("id", "", "")
	tf := f.String("type", "tcp", "")
	timeoutf := f.Duration("timeout", 0, "")
	desc := f.String("desc", "upnpctl v"+VERSION, "")
	//parse and transform args
	f.Parse(args)

	args = f.Args()

	timeout := int((*timeoutf).Seconds())
	t := upnp.Protocol(strings.ToUpper(*tf))
	switch t {
	case upnp.TCP:
	case upnp.UDP:
	default:
		display("Invalid type: " + string(t))
	}

	l := len(args)
	plural := "s"
	if l == 1 {
		plural = ""
	}

	ms := make([]*mapping, l)
	for i, a := range args {
		m := &mapping{}
		if err := m.unmarshal(a); err != nil {
			display(err.Error())
		}
		if cmd == rem && m.internal != m.external {
			display("When removing ports, only specify the external port")
		}
		// fmt.Printf("Mapping %d -> %d (timeout %d, description %s)\n", m.external, m.internal, timeout, *desc)
		ms[i] = m
	}

	var c *client = nil
	fmt.Printf("Discovering UPnP devices...\n")
	cs := discover()
	if len(cs) == 0 {
		display("No UPnP devices found")
	}

	if *id == "" {
		if len(cs) == 1 {
			c = cs[0]
		} else {
			fmt.Printf("The --id option is required as there is more than one UPnP device:\n")
			for _, c := range cs {
				fmt.Printf("  --id %s => %s (%s)\n", c.id, c.name, c.ip)
			}
			os.Exit(1)
		}
	} else {
		for _, cl := range cs {
			if cl.id == *id {
				c = cl
				break
			}
		}
		if c == nil {
			display("No UPnP devices found matching id: " + *id)
		}
	}

	if cmd == add {
		fmt.Printf("Adding #%d mapping%s...\n", l, plural)
		for _, m := range ms {
			err := c.igd.AddPortMapping(t, m.external, m.internal, *desc, timeout)
			if err != nil {
				display(fmt.Sprintf("Failed to add mapping %d:%d (%s)", m.external, m.internal, err))
			}
		}
	}

	if cmd == rem {
		fmt.Printf("Removing #%d mapping%s...\n", l, plural)
		for _, m := range ms {
			err := c.igd.DeletePortMapping(t, m.external)
			if err != nil {
				display(fmt.Sprintf("Failed to remove mapping %d (%s)", m.external, err))
			}
		}
	}

	fmt.Println("Done")
}

func listMappings() {
	fmt.Printf("Discovering UPnP devices...\n")
	cs := discover()
	for _, c := range cs {
		fmt.Printf("  #%s: %s (%s)\n", c.id, c.name, c.ip)
	}
}

func display(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func discover() clients {
	cs := make(clients, 0)
	igds := upnp.Discover()
	for _, igd := range igds {
		host := igd.URL().Host
		ip, _, _ := net.SplitHostPort(host)
		h := sha1.New()
		h.Write([]byte(igd.UUID()))
		h.Write([]byte(host))
		hex := hex.EncodeToString(h.Sum(nil))
		id := hex[0:5] //only use first 5 chars (of 40)
		cs = append(cs, &client{igd, igd.FriendlyName(), ip, id})
	}
	return cs
}

type mapping struct {
	external int
	internal int
}

func (m *mapping) unmarshal(s string) error {
	ports := strings.SplitN(s, ":", 2)
	var err error
	if len(ports) == 1 {
		m.external, err = strconv.Atoi(ports[0])
		if err != nil || !valid(m.external) {
			return fmt.Errorf("Invalid port '%s'", ports[0])
		}
		m.internal = m.external
	} else {
		m.external, err = strconv.Atoi(ports[0])
		if err != nil || !valid(m.external) {
			return fmt.Errorf("Invalid external port '%s'", ports[0])
		}
		m.internal, err = strconv.Atoi(ports[1])
		if err != nil || !valid(m.internal) {
			return fmt.Errorf("Invalid internal port '%s'", ports[1])
		}
	}
	return nil
}

type clients []*client

type client struct {
	igd          upnp.IGD
	name, ip, id string
}

func valid(port int) bool {
	return port > 0 && port < 65536
}
