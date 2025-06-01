package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func LogRequest(c *fiber.Ctx) error {
	jsonHeaders := make(map[string]string)

	c.Request().Header.VisitAll(func(key, value []byte) {
		jsonHeaders[string(key)] = string(value)
	})

	logger.Log.Debug("Request Headers", "headers", jsonHeaders)

	return c.Next()
}
