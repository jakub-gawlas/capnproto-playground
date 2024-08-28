# Cap'n'proto playground

## Installation

```bash
brew install capnp
go install capnproto.org/go/capnp/v3/capnpc-go@latest
go get capnproto.org/go/capnp/v3
```

## Development

### Compile schema

```bash
capnp compile -I (go list -m -f '{{.Dir}}' capnproto.org/go/capnp/v3)/std -ogo src/book/book.capnp
```
