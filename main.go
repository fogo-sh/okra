package main

import (
	"log"

	"github.com/fogo-sh/okra/app"
	"github.com/fogo-sh/okra/app/database"
	"github.com/fogo-sh/okra/app/routes"
)

func main() {
	okraApp, err := app.Setup()
	if err != nil {
		log.Fatalf("Error configuring app: %s", err.Error())
		return
	}

	routes.Register(okraApp.Fiber)

	err = database.CreateInstance(okraApp.Config.DBPath)
	if err != nil {
		log.Fatalf("Error creating database instance: %s", err.Error())
		return
	}

	app.Start()
}
