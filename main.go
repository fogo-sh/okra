package main

import (
	"fmt"

	"github.com/fogo-sh/okra/app"
	"github.com/fogo-sh/okra/app/database"
	"github.com/fogo-sh/okra/app/routes"
	"github.com/rs/zerolog/log"
)

// CreateTestData initializes the database with testing data.
func CreateTestData() error {
	log.Info().Msg("Creating testing data")
	user, err := database.NewUser("test", "test")
	if err != nil {
		return fmt.Errorf("error creating test user: %w", err)
	}
	podcast, err := database.NewPodcast(user.ID, "https://feeds.megaphone.fm/darknetdiaries")
	if err != nil {
		return fmt.Errorf("error creating test podcast: %w", err)
	}
	err = podcast.Refresh()
	if err != nil {
		return fmt.Errorf("error refreshing test podcast: %w", err)
	}
	log.Info().Msg("Test data created")
	return nil
}

func main() {
	okraApp, err := app.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Error configuring app")
		return
	}

	routes.Register(okraApp.Fiber)

	err = database.CreateInstance(okraApp.Config.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating database instance")
		return
	}
	database.Migrate()

	if okraApp.Config.Debug {
		err = CreateTestData()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to create test data")
		}
	}

	app.Start()
}
