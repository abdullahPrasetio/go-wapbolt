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

	// Register Metadata for Login (Validations & Examples)
	wapbolt.RegisterMetadata("POST", "/api/v1/auth/login", wapbolt.RouteMetadata{
		Description: "Authenticate user and get access token",
		FieldValidations: map[string]interface{}{
			"body": map[string]interface{}{
				"email": wapbolt.ValidationRule{
					Rules:       []string{"required", "email"},
					Description: "Registered user email address",
				},
				"password": wapbolt.ValidationRule{
					Rules:       []string{"required"},
					Description: "User password (min 8 chars)",
				},
			},
		},
		Examples: []wapbolt.Response{
			{
				Name:   "Success Response",
				Status: "OK",
				Code:   200,
				Body:   "{\n  \"token\": \"eyJhbGci...\",\n  \"expires_in\": 3600\n}",
			},
			{
				Name:   "Invalid Credentials",
				Status: "Unauthorized",
				Code:   401,
				Body:   "{\n  \"error\": \"invalid_credentials\"\n}",
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
