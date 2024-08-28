package server

import (
	"capnproto-playground/src/calculator"
	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"context"
	"log/slog"
	"net"
)

func Run(ctx context.Context) error {
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		return err
	}

	for {
		nc, err := l.Accept()
		if err != nil {
			return err
		}

		go func(nc net.Conn) {
			log := slog.With(slog.Any("remote_addr", nc.RemoteAddr()))
			log.Info("connection accepted")

			c := calculator.Calculator_ServerToClient(calculator.Server{})
			conn := rpc.NewConn(rpc.NewStreamTransport(nc), &rpc.Options{
				BootstrapClient: capnp.Client(c),
			})
			defer conn.Close()

			<-conn.Done()
			log.Info("connection closed")
		}(nc)
	}
}
