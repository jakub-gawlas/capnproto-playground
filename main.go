package main

import (
	"capnproto-playground/src/cmd/client"
	"capnproto-playground/src/cmd/server"
	"context"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing command")
		return
	}
	cmd := os.Args[1]
	ctx := context.Background()
	switch cmd {
	case "server":
		if err := server.Run(ctx); err != nil {
			panic(err)
		}
	case "client":
		if err := client.Run(ctx); err != nil {
			panic(err)
		}
	default:
		fmt.Println("unknown command")
	}
}
