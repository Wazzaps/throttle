# Throttle

A simple utility to simulate latency in UNIX domain sockets.

Used for testing high latency links for [BoldUI](https://github.com/Wazzaps/boldui).

### Usage example:

- `go build`
- `./throttle --listen=/tmp/server-slow.sock --connect=/tmp/server.sock --latency=1000`
- Terminal #1: `nc -vlU /tmp/server.sock`
- Terminal #2: `nc -vU /tmp/server-slow.sock`

The data is delayed in both directions now, by 1000ms.

