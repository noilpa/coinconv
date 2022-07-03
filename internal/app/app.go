package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"coinconv/internal/client/coinmarketcap"
	"coinconv/internal/conventer"
)

type iConvert interface {
	Convert(ctx context.Context, num float64, from, to string) (float64, error)
}

type app struct {
	c iConvert
}

func New() (*app, error) {
	c, err := coinmarketcap.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create coinmarketcap client: %v\n", err)
	}

	return &app{
		c: conventer.New(c),
	}, nil
}

func (a *app) Run() int {
	in, err := parseArgs(os.Args[1:])
	if err != nil {
		printHelp()
		log.Fatalf("failed to parse args: %v", err)
	}

	ctx := context.Background()

	res, err := a.c.Convert(ctx, in.num, in.from, in.to)
	if err != nil {
		log.Fatalf("failed to convert: %v\n", err)
	}

	fmt.Printf("%.4f %s -> %.4f %s\n", in.num, in.from, res, in.to)

	return 0
}

type input struct {
	from string
	to   string
	num  float64
}

func parseArgs(args []string) (in input, err error) {
	if len(args) != 3 {
		return in, fmt.Errorf("wrong args number: %d != 3", len(args))
	}

	in.num, err = strconv.ParseFloat(args[0], 64)
	if err != nil {

		return in, fmt.Errorf("failed to parse currency number: %v\n", err)
	}
	in.from = args[1]
	in.to = args[2]

	return in, nil
}

func printHelp() {
	fmt.Println("example of usage: ./coinconv 123.45 USD BTC")
}
