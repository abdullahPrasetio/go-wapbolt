package main

import (
	"log"

	"github.com/abdullahPrasetio/go-wapbolt/adapters/gin_adapter"
	"github.com/abdullahPrasetio/go-wapbolt/examples/gin_zero_touch/app/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 1. Group rute
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", controllers.RegisterUser)
	}

	// 2. OTOMATIS TAMBAHKAN ROUTE DOKUMENTASI (Cuma 1 baris!)
	gin_adapter.Register(r, "/export-docs", "http://localhost:8080")

	log.Fatal(r.Run(":8080"))
}
