package calculator

import (
	"capnproto.org/go/capnp/v3"
	"context"
	"fmt"
	"log/slog"
)

type Server struct{}

func (s Server) Add(ctx context.Context, call Calculator_add) error {
	return s.op(ctx, call, "add", func(a, b int32) (int32, error) {
		return a + b, nil
	})
}

func (s Server) Sub(ctx context.Context, call Calculator_sub) error {
	return s.op(ctx, call, "sub", func(a, b int32) (int32, error) {
		return a - b, nil
	})
}

func (s Server) Mul(ctx context.Context, call Calculator_mul) error {
	return s.op(ctx, call, "mul", func(a, b int32) (int32, error) {
		return a * b, nil
	})
}

func (s Server) Div(ctx context.Context, call Calculator_div) error {
	return s.op(ctx, call, "div", func(a, b int32) (int32, error) {
		if b == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return a / b, nil
	})
}

type opCall interface {
	Args() Calculator_InputPair
	AllocResults() (Calculator_Output, error)
}

func (s Server) op(ctx context.Context, call opCall, opName string, opFn func(a, b int32) (int32, error)) error {
	left, right, err := resolveInputPair(ctx, call.Args())
	if err != nil {
		return err
	}

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	result, err := opFn(left, right)
	if err != nil {
		return err
	}

	slog.Info(opName, slog.Any("left", left), slog.Any("right", right), slog.Any("result", result))

	if err := res.SetValue(newValue(result)); err != nil {
		return err
	}

	return nil
}

type valueServer struct {
	value int32
}

func (v valueServer) Read(ctx context.Context, call Calculator_Value_read) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	res.SetValue(v.value)
	return nil
}

func NewLiteralInput(literal int32) (Calculator_Input, error) {
	arena := capnp.SingleSegment(nil)

	_, seg, err := capnp.NewMessage(arena)
	if err != nil {
		panic(err)
	}

	r, err := NewRootCalculator_Input(seg)
	if err != nil {
		panic(err)
	}

	r.SetLiteral(literal)

	return r, nil
}

func NewValueInput(value Calculator_Value) (Calculator_Input, error) {
	arena := capnp.SingleSegment(nil)

	_, seg, err := capnp.NewMessage(arena)
	if err != nil {
		panic(err)
	}

	r, err := NewRootCalculator_Input(seg)
	if err != nil {
		panic(err)
	}

	if err := r.SetValue(value); err != nil {
		return Calculator_Input{}, err
	}

	return r, nil
}

func newValue(value int32) Calculator_Value {
	return Calculator_Value_ServerToClient(valueServer{value})
}

func resolveInputPair(ctx context.Context, inputPair Calculator_InputPair) (int32, int32, error) {
	left, err := inputPair.Left()
	if err != nil {
		return 0, 0, err
	}

	leftValue, err := resolveInput(ctx, left)
	if err != nil {
		return 0, 0, err
	}

	right, err := inputPair.Right()
	if err != nil {
		return 0, 0, err
	}

	rightValue, err := resolveInput(ctx, right)
	if err != nil {
		return 0, 0, err
	}

	return leftValue, rightValue, nil
}

func resolveInput(ctx context.Context, input Calculator_Input) (int32, error) {
	switch input.Which() {
	case Calculator_Input_Which_literal:
		return input.Literal(), nil
	case Calculator_Input_Which_value:
		f, release := input.Value().Read(ctx, nil)
		defer release()
		r, err := f.Struct()
		if err != nil {
			return 0, err
		}
		return r.Value(), nil
	}
	return 0, fmt.Errorf("empty input")
}
