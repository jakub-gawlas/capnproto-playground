# Cap'n'proto playground

This repo demonstrates a client-server application using [Cap'n'Proto](https://capnproto.org/) protocol.

Seems trivial, but I have spent a bit of time to figure out how to implement and run it properly. Cap'n'Proto documentation wasn't very helpful.

## Usage

The example implements a simple calculator. The server provides the calculator service, and the client executes calculation requests and fetches the result at the end. 

Start the server:

```bash
go run main.go server
```

Start the client which run a calculation `(90 - (179 + 233) / 7) * 3`:

```bash
go run main.go client
```

## Development

### Installation

```bash
brew install capnp
go install capnproto.org/go/capnp/v3/capnpc-go@latest
go get capnproto.org/go/capnp/v3
```

### Compile schema

Example for book package:

```bash
capnp compile -I (go list -m -f '{{.Dir}}' capnproto.org/go/capnp/v3)/std -ogo src/book/book.capnp
```