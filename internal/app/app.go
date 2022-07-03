package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"coinconv/internal/client/coinmarketcap"
	"coinconv/internal/config"
	"coinconv/internal/conventer"
)

type app struct {
	cfg config.Config
}

func New(cfg config.Config) *app {
	return &app{
		cfg: cfg,
	}
}

func (a *app) Run() int {
	args := os.Args[1:]

	if len(args) != 3 {
		printHelp()
		return 1
	}

	num, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		printHelp()
		log.Fatalf("failed to parse currency number: %v\n", err)
	}
	from := args[1]
	to := args[2]

	c, err := coinmarketcap.New(a.cfg.CoinmarketcapHost, a.cfg.CoinmarketcapSecretKey)
	if err != nil {
		log.Fatalf("failed to create coinmarketcap client: %v", err)
	}

	conv := conventer.New(c)

	ctx := context.Background()

	res, err := conv.Convert(ctx, num, from, to)
	if err != nil {
		log.Fatalf("failed to convert: %v", err)
	}

	fmt.Printf("%.4f %s -> %.4f %s\n", num, from, res, to)

	return 0
}

func printHelp() {
	fmt.Println("example of usage: ./coinconv 123.45 USD BTC")
}
