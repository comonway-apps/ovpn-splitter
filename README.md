# OpenVPN splitter

This CLI utility split opvn file from OpenVPN Access Server to several files:

- readme
- ca, 
- cert, 
- key
- tls-auth

**Just do**

``` shell
$ go build
$ ovpn-splitter path-to-opvn.ovpn
```

If you launch app with no args, OpenVPN splitter search for `client.ovpn` file.
