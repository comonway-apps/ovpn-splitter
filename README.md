# OpenVPN splitter

This CLI utility split opvn file from OpenVPN Access Server to several files:

- readme
- ca, 
- cert, 
- key
- tls-auth

## Using OpenVPN splitter

**For Windows**
Drag-and-Drop your opvn file on ovpn-splitter.exe

or

``` shell
$ go build
$ ovpn-splitter path-to-opvn.ovpn
```

If you launch app with no args, OpenVPN splitter search for `client.ovpn` file.
