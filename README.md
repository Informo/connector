# Informo connector

The Informo connector is a basic tool, usable as a single binary, which will connect a Matrix homeserver to the Informo network.

It will:

* Create a user on the homeserver (the `m.login.dummy` auth flow must be enabled)
* Make it join `#informo:matrix.org`
* Make it create an entry point to the network (by appending it to the list of entrypoints for this server)

Of course, all this steps can be done by hand. The connector is just a program that will run them for you.

To use it, run:

```
informo-connector --server-name matrix.example.tld
```

All command line options are optional. If not provided, the server name defaults to `127.0.0.1`.

The connector will perform a [DNS lookup on the federation record](https://github.com/matrix-org/synapse#setting-up-federation) to find the FQDN and port at which the Matrix homeserver is reachable. If it doesn't find any record, it will fallback on the value provided as the server name for the FQDN, and `8448` for the port.

These values can be overriden by providing either the server's FQDN (with the `--fqdn` command line option), its port (with the `--port` command line option) or both of them. If only one of these values is provided the connector will perform a DNS lookup to try to fill the unprovided value (falling back as previously stated if the lookup's result is empty), but if both of them are provided the DNS lookup won't be performed at all.

The connector will try to use HTTPS by default when communicating with the homeserver. If you want to disable TLS and use plain text HTTP instead, append the `--no-tls` option to your call. **Warning: this is not recommended in other cases than local development! If you disable TLS, traffic will be sent in plain text (not encrypted) and anyone spying on your network will be able to disconnect your server from the Informo network!**

## Build

Because of Go's poor dependencies and workspaces management, the connector uses `gb`. You can install `gb` by running:

```
go get github.com/constabulary/gb/...
```

You can then compile the connector by running `gb build` at the root of this repository (where this README is located).

If this doesn't work, make sure your `$PATH` contains your `$GOPATH`.
