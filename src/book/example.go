package book

import (
	"capnproto.org/go/capnp/v3"
	"fmt"
)

func main() {
	arena := capnp.SingleSegment(nil)

	msg, seg, err := capnp.NewMessage(arena)
	if err != nil {
		panic(err)
	}

	b, err := NewRootBook(seg)
	if err != nil {
		panic(err)
	}

	if err := b.SetTitle("The Go Programming Language"); err != nil {
		panic(err)
	}

	b.SetPages(380)

	data, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	fmt.Println("len:", len(data))

	data, err = msg.MarshalPacked()
	if err != nil {
		panic(err)
	}
	fmt.Println("len packed:", len(data))

	msg2, err := capnp.UnmarshalPacked(data)
	if err != nil {
		panic(err)
	}

	b2, err := ReadRootBook(msg2)
	if err != nil {
		panic(err)
	}

	b.SetPages(100)

	fmt.Println(b.Pages() == b2.Pages())
}
