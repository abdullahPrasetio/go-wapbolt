package main

import (
	"log"

	"github.com/abdullahPrasetio/go-wapbolt/adapters/fiber_adapter"
	v1 "github.com/abdullahPrasetio/go-wapbolt/examples/zero_touch_example/controllers/v1"
	v2 "github.com/abdullahPrasetio/go-wapbolt/examples/zero_touch_example/controllers/v2"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// 1. API V1 & V2 Groups
	apiV1 := app.Group("/api/v1")
	apiV1.Post("/users", v1.CreateUser)
	apiV1.Put("/users", v1.UpdateUser)

	apiV2 := app.Group("/api/v2")
	apiV2.Post("/users", v2.CreateUser)

	// 2. OTOMATIS TAMBAHKAN ROUTE DOKUMENTASI (Cuma 1 baris!)
	fiber_adapter.Register(app, "/export-docs", "http://localhost:3000")

	log.Fatal(app.Listen(":3000"))
}
