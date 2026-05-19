package v1

import (
	"github.com/abdullahPrasetio/go-wapbolt/examples/zero_touch_example/models"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	return c.JSON(req)
}

func UpdateUser(c *fiber.Ctx) error {
	var req models.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	return c.JSON(req)
}
