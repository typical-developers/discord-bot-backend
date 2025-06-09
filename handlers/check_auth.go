package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func CheckAuthorization(c *fiber.Ctx) error {
	auth := c.Get("X-API-KEY")

	logger.Log.Debug("Authorizing Request", "auth", auth, "key", os.Getenv("API_KEY"))
	if auth != os.Getenv("API_KEY") {
		logger.Log.Warn("An unauthorized request was made.", "ip", c.IP())

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	return c.Next()
}
