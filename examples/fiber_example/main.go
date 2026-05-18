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

	// Define a struct for Login request
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email" description:"User login email"`
		Password string `json:"password" validate:"required,min=8" description:"Secret password"`
	}

	auth := api.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error { return c.SendString("Login") })

	// Register Metadata for Login using AUTOMATIC Struct Parsing
	wapbolt.RegisterMetadata("POST", "/api/v1/auth/login", wapbolt.RouteMetadata{
		Description: "Authenticate user using struct tags",
		FieldValidations: map[string]interface{}{
			"body": wapbolt.ParseStruct(LoginRequest{}),
		},
		Examples: []wapbolt.Response{
			{
				Name:   "Success",
				Status: "OK",
				Code:   200,
				Body:   "{\n  \"token\": \"abc-123\"\n}",
			},
		},
	})

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
