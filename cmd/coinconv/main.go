package main

import (
	"log"
	"os"

	"coinconv/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %v\n", err)
	}
	os.Exit(a.Run())
}
