package routes

import "github.com/gofiber/fiber/v2"

func _IndexHandler(c *fiber.Ctx) error {
	return c.Render("index", &fiber.Map{})
}

// Register all of the app's routes with Fiber.
func Register(app *fiber.App) {
	app.Get("/", _IndexHandler)
}
