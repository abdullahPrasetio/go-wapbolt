package main

import (
	"log"

	"github.com/abdullahPrasetio/go-wapbolt/adapters/gin_adapter"
	"github.com/abdullahPrasetio/go-wapbolt/wapbolt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Define some routes
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.GET("/", func(c *gin.Context) { c.String(200, "Get all users") })
			users.POST("/", func(c *gin.Context) { c.String(200, "Create user") })
			users.GET("/:id", func(c *gin.Context) { c.String(200, "Get user by ID") })
		}
		
		auth := api.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) { c.String(200, "Login") })
		}
	}

	// Endpoint to trigger Wapbolt Collection generation
	r.GET("/generate-wapbolt", func(c *gin.Context) {
		col := wapbolt.NewCollection("Gin Example API", "Auto-generated collection from Gin")
		
		// Parse routes into the collection
		gin_adapter.ParseRoutes(r, col, "http://localhost:8080")

		// Export to file
		filename := "wapbolt-gin-collection.json"
		err := col.ExportToFile(filename)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.String(200, "Collection generated: "+filename)
	})

	log.Fatal(r.Run(":8080"))
}
