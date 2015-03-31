# upnpctl

A small UPnP client

:warning: Beta

### Install

**Binaries**

See [latest release](https://github.com/jpillora/upnpctl/releases/latest)

**Source**

``` sh
$ go get -v github.com/jpillora/upnpctl
```

### Examples

Forward the router's port 3000 to this machine's port 3000

```
upnpctl add 3000
```

Forward the router's port 4000 to this machine's port 5000

```
upnpctl add 4000:5000
```

### Usage

```
$ upnpctl --help
```

<tmpl,code: go run main.go --help>
```

	Usage: upnpctl <command> [options]
	
	Version: 0.0.0-src

	Commands:
	  * list: discovers all available UPnP devices
	  * add: adds a set of port mappings to a device
	  * rem: removes a set of port mappings from a device

	Options:

	  -v, verbose logs
	  -vv, very verbose logs

	Read more: https://github.com/jpillora/upnpctl

```
</tmpl>

#### MIT License

Copyright © 2015 Jaime Pillora &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.