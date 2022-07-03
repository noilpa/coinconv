package main

import (
	"log"
	"os"

	"coinconv/internal/app"
	"coinconv/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(app.New(cfg).Run())
}
