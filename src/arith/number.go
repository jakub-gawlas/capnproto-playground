package arith

import (
	"capnproto.org/go/capnp/v3"
	"context"
	"log/slog"
	"math/rand"
)

type NumberOpsServer struct {
	number Number
}

func NewNumberOpsServer() (*NumberOpsServer, error) {
	value := rand.Int63n(1e3)
	number, err := newNumber(value)
	if err != nil {
		return nil, err
	}

	return &NumberOpsServer{
		number: number,
	}, nil
}

func (n *NumberOpsServer) Plus(ctx context.Context, call Number_Ops_plus) error {
	value := n.number.Value()
	valueToAdd := call.Args().ValueToAdd()
	sum := value + valueToAdd

	slog.Info("called: sum on number", slog.Int64("number_value", value), slog.Int64("value_to_add", valueToAdd), slog.Int64("sum", sum))

	n.number.SetValue(sum)

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	if err := res.SetSum(n.number); err != nil {
		return err
	}

	return nil
}

func newNumber(value int64) (Number, error) {
	arena := capnp.SingleSegment(nil)

	_, seg, err := capnp.NewMessage(arena)
	if err != nil {
		return Number{}, err
	}

	b, err := NewRootNumber(seg)
	if err != nil {
		return Number{}, err
	}

	b.SetValue(value)

	return b, nil
}
