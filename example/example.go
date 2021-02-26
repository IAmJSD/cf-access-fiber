package main

import (
	"github.com/gofiber/fiber/v2"
	access "github.com/jakemakesstuff/cf-access-fiber"
	"log"
	"os"
)

// Initialises the middleware.
var middleware = access.Validate(
	os.Getenv("TEAM_DOMAIN"), os.Getenv("APPLICATION_AUD"), unauthorized)

// The route that is hit when a unauthenticated/invalid Access request is made.
func unauthorized(ctx *fiber.Ctx) error {
	ctx.Status(400)
	_, err := ctx.WriteString("No content for you!")
	return err
}

// The main application.
func main() {
	// Create the Fiber application.
	app := fiber.New()

	// Insert the middleware. You could do this per route too.
	app.Use(middleware)

	// Add a basic route to /.
	app.Get("/", func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(*access.CloudflareAccessUserInfo)
		_, err := ctx.WriteString(
			"You are " + user.Email + " (user ID: " + user.UserIdentifier + ")")
		return err
	})

	// Start the server.
	log.Fatal(app.Listen(":8080"))
}
