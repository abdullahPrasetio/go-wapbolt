# go-wapbolt

`go-wapbolt` is an unofficial Golang SDK / Module designed to automatically generate API Collections for Wapbolt from popular Go web frameworks. It extracts your defined routes and exports them into a `.json` collection file that can be easily imported into the Wapbolt application.

## Features
- **Core Collection Builder**: Create Wapbolt collections, folders, and requests programmatically.
- **Fiber Adapter**: Automatically parse all routes from a `*fiber.App` instance.
- **Gin Adapter**: Automatically parse all routes from a `*gin.Engine` instance.

## Installation

```bash
go get github.com/abdullahPrasetio/go-wapbolt
```

## Quick Start (Fiber)

```go
package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/abdullahPrasetio/go-wapbolt/adapters/fiber_adapter"
	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
)

func main() {
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error { return c.SendString("Get Users") })
	app.Post("/users", func(c *fiber.Ctx) error { return c.SendString("Create User") })

	// Endpoint to generate Wapbolt collection
	app.Get("/generate-docs", func(c *fiber.Ctx) error {
		col := wapbolt.NewCollection("My API", "Generated from Fiber")
		
		// 1. Parse routes
		fiber_adapter.ParseRoutes(app, col, "http://localhost:3000")
		
		// 2. Export to JSON file
		err := col.ExportToFile("wapbolt-collection.json")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		
		return c.SendString("Collection generated successfully!")
	})

	log.Fatal(app.Listen(":3000"))
}
```

## Quick Start (Gin)

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/abdullahPrasetio/go-wapbolt/adapters/gin_adapter"
	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
)

func main() {
	r := gin.Default()

	r.GET("/items", func(c *gin.Context) {})
	r.POST("/items", func(c *gin.Context) {})

	// Endpoint to generate Wapbolt collection
	r.GET("/generate-docs", func(c *gin.Context) {
		col := wapbolt.NewCollection("My Gin API", "Generated from Gin")
		
		// 1. Parse routes
		gin_adapter.ParseRoutes(r, col, "http://localhost:8080")
		
		// 2. Export to JSON file
		err := col.ExportToFile("wapbolt-collection.json")
		if err != nil {
			c.String(500, err.Error())
			return
		}
		
		c.String(200, "Collection generated successfully!")
	})

	log.Fatal(r.Run(":8080"))
}
```

## License
MIT
