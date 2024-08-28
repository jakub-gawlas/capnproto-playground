package client

import (
	"capnproto-playground/src/calculator"
	"capnproto.org/go/capnp/v3/rpc"
	"context"
	"fmt"
	"net"
)

func Run(ctx context.Context) error {
	nc, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		return err
	}

	conn := rpc.NewConn(rpc.NewStreamTransport(nc), nil)
	defer conn.Close()

	calc := calculator.Calculator(conn.Bootstrap(ctx))

	// (90 - (179 + 233) / 7) * 3

	a, release := calc.Add(ctx, inputPair(withLeftLiteral(179), withRightLiteral(233)))
	defer release()

	b, release := calc.Div(ctx, inputPair(withLeftValue(a.Value()), withRightLiteral(7)))
	defer release()

	c, release := calc.Sub(ctx, inputPair(withLeftLiteral(90), withRightValue(b.Value())))
	defer release()

	d, release := calc.Mul(ctx, inputPair(withLeftValue(c.Value()), withRightLiteral(3)))
	defer release()

	result, release := d.Value().Read(ctx, nil)
	defer release()

	rs, err := result.Struct()
	if err != nil {
		return err
	}

	fmt.Println(rs.Value())

	return nil
}

type options struct {
	leftLiteral  *int32
	leftValue    *calculator.Calculator_Value
	rightLiteral *int32
	rightValue   *calculator.Calculator_Value
}

func withLeftLiteral(literal int32) func(*options) {
	return func(o *options) {
		o.leftLiteral = &literal
	}
}

func withLeftValue(value calculator.Calculator_Value) func(*options) {
	return func(o *options) {
		o.leftValue = &value
	}
}

func withRightLiteral(literal int32) func(*options) {
	return func(o *options) {
		o.rightLiteral = &literal
	}
}

func withRightValue(value calculator.Calculator_Value) func(*options) {
	return func(o *options) {
		o.rightValue = &value
	}
}

func inputPair(opts ...func(*options)) func(calculator.Calculator_InputPair) error {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return func(p calculator.Calculator_InputPair) error {
		var err error

		var leftInput calculator.Calculator_Input
		if o.leftLiteral != nil {
			leftInput, err = calculator.NewLiteralInput(*o.leftLiteral)
			if err != nil {
				return err
			}
		} else if o.leftValue != nil {
			leftInput, err = calculator.NewValueInput(*o.leftValue)
			if err != nil {
				return err
			}
		}
		if err := p.SetLeft(leftInput); err != nil {
			return err
		}

		var rightInput calculator.Calculator_Input
		if o.rightLiteral != nil {
			rightInput, err = calculator.NewLiteralInput(*o.rightLiteral)
			if err != nil {
				return err
			}
		} else if o.rightValue != nil {
			rightInput, err = calculator.NewValueInput(*o.rightValue)
			if err != nil {
				return err
			}
		}
		if err := p.SetRight(rightInput); err != nil {
			return err
		}

		return nil
	}
}
