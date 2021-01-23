package main

import (
	"log"

	"github.com/fogo-sh/okra/app"
)

func main() {
	err := app.Setup()
	if err != nil {
		log.Fatalf("Error configuring app: %s", err.Error())
	}
	app.Start()
}
