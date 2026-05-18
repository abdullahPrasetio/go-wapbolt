package main

import (
	"log"

	"github.com/abdullahPrasetio/go-wapbolt/adapters/fiber_adapter"
	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Define some routes
	api := app.Group("/api/v1")
	
	users := api.Group("/users")
	users.Get("/", func(c *fiber.Ctx) error { return c.SendString("Get all users") })
	users.Post("/", func(c *fiber.Ctx) error { return c.SendString("Create user") })
	users.Get("/:id", func(c *fiber.Ctx) error { return c.SendString("Get user by ID") })

	auth := api.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error { return c.SendString("Login") })

	// Endpoint to trigger Wapbolt Collection generation
	app.Get("/generate-wapbolt", func(c *fiber.Ctx) error {
		col := wapbolt.NewCollection("Fiber Example API", "Auto-generated collection from Fiber")
		
		// Parse routes into the collection
		fiber_adapter.ParseRoutes(app, col, "http://localhost:3000")

		// Export to file
		filename := "wapbolt-fiber-collection.json"
		err := col.ExportToFile(filename)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.SendString("Collection generated: " + filename)
	})

	log.Fatal(app.Listen(":3000"))
}
