package arith

import (
	"context"
	"errors"
	"log/slog"
)

type Server struct{}

func (Server) Multiply(ctx context.Context, call Arith_multiply) error {
	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	a := call.Args().A()
	b := call.Args().B()
	product := a * b

	res.SetProduct(product)
	slog.Info("called: multiply", slog.Int64("a", a), slog.Int64("b", b), slog.Int64("product", product))

	return nil
}

func (Server) Divide(ctx context.Context, call Arith_divide) error {
	if call.Args().Denom() == 0 {
		return errors.New("divide by zero")
	}

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	res.SetQuotient(call.Args().Num() / call.Args().Denom())
	res.SetRemainder(call.Args().Num() % call.Args().Denom())
	return nil
}
