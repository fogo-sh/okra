package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/markbates/pkger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CurrentApp represents the current app instance.
var CurrentApp *App

// Config represents the basic configuration required for the app on startup.
type Config struct {
	ListenAddress string `default:":3000" split_words:"true"`
	DBPath        string `default:"file::memory:?cache=shared" split_words:"true"`
	Debug         bool   `default:"false"`
}

// App represents an app instance and all of it's required resources.
type App struct {
	Fiber  *fiber.App
	DB     *gorm.DB
	Config Config
}

// Setup sets up the global app instance.
func Setup() error {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load .env file: %s\n", err.Error())
	}

	var config Config
	err = envconfig.Process("okra", &config)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening db connection: %w", err)
	}

	engine := html.NewFileSystem(pkger.Dir("/app/templates"), ".html")
	engine.Reload(config.Debug)

	fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true, Views: engine})

	CurrentApp = &App{fiberApp, db, config}

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", &fiber.Map{})
	})

	return nil
}

// Start starts the app. Requires calling of Setup first.
func Start() {
	if CurrentApp == nil {
		panic(errors.New("unable to start uninitialized app - call app.Setup first"))
	}

	log.Info().Msg("Starting Okra")
	CurrentApp.Fiber.Listen(CurrentApp.Config.ListenAddress)
}