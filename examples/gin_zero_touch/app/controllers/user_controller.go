package controllers

import (
	"github.com/abdullahPrasetio/go-wapbolt/examples/gin_zero_touch/app/models"
	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var req models.RegisterRequest
	// Scanner akan mendeteksi BindJSON dan melacak models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, req)
}
