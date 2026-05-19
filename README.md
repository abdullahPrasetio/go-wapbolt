# go-wapbolt

`go-wapbolt` is a Golang SDK designed to automatically generate API Collections for Wapbolt using **Zero-Touch Static Analysis**. It scans your source code (Fiber or Gin), detects routes, resolves request models, and parses validation tags automatically.

## Features
- **Zero-Touch Generation**: No need to modify your business logic. Just scan and export.
- **Smart AST Scanner**: Automatically detects routes and matches them with `BodyParser` or `BindJSON` calls.
- **Validation Tag Support**: Reads `json`, `validate`, and `description` tags from your structs.
- **Multi-Package Support**: Works across nested folders (Controllers/Models).
- **Auto Sample Body**: Generates sample JSON payloads based on your validation rules.
- **Framework Adapters**: Full support for **Fiber** and **Gin**.

## Installation

```bash
go get github.com/abdullahPrasetio/go-wapbolt
```

## Quick Start (Zero-Touch)

You only need to call `AutoScan` before exporting. It will recursively scan your project directory.

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
	"github.com/abdullahPrasetio/go-wapbolt/adapters/fiber_adapter"
)

func main() {
	app := fiber.New()

	// 1. Your normal routes
	app.Post("/api/users", createUserHandler)

	// 2. Add an export endpoint
	app.Get("/export-docs", func(c *fiber.Ctx) error {
		// Scans the entire project folder
		wapbolt.AutoScan(".") 

		col := wapbolt.NewCollection("My API", "v1.0")
		fiber_adapter.ParseRoutes(app, col, "http://localhost:3000")
		col.ExportToFile("wapbolt-collection.json")
		
		return c.SendString("Collection generated!")
	})

	app.Listen(":3000")
}
```

## How it Works
The scanner uses Go's `Abstract Syntax Tree (AST)` to:
1. Find all struct definitions and their tags.
2. Find all function handlers and look for body parsing calls like `c.BodyParser(&MyRequest{})`.
3. Associate the route path with the detected struct.
4. Normalize everything into a Wapbolt-compatible JSON collection.

## License
MIT
