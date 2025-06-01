package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func CheckAuthorization(c *fiber.Ctx) error {
	// auth := c.Get("X-API-KEY")
	// if auth != os.Getenv("API_KEY") {
	// 	return c.SendStatus(fiber.StatusUnauthorized)
	// }

	return c.Next()
}
