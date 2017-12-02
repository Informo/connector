# Informo connector

The Informo connector is a basic tool, usable as a single binary, which will connect a Matrix homeserver to the Informo network.

It will:

* Create a user on the homeserver (the `m.login.dummy` auth flow must be enabled)
* Make it join `#informo:matrix.org`
* Make it create an entry point to the network (by appending it to the list of entrypoints for this server)

Of course, all this steps can be done by hand. The connector is just a program that will run them for you.

To use it, run:

```
informo-connector --homeserver matrix.example.tld --port 8448
```

All command line options are optional and, if not provided, default to `127.0.0.1` for the homeserver and `443` for the port.

The connector will try to use HTTPS by default. If you want to disable TLS and use plain text HTTP instead, append the `--no-tls` option to your call. **Warning: this is not recommended in other cases than local development! If you disable TLS, traffic will be sent in plain text (not encrypted) and anyone spying on your network will be able to disconnect your server from the Informo network!**

## Build

Because of Go's poor dependencies and workspaces management, the connector uses `gb`. You can install `gb` by running:

```
go get github.com/constabulary/gb/...
```

You can then compile the connector by running `gb build` at the root of this repository (where this README is located).

If this doesn't work, make sure your `$PATH` contains your `$GOPATH`.
