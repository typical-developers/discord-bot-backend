package handlers

import (
	"github.com/gofiber/fiber/v2"
	api_structures "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func SetPoolConn(c *fiber.Ctx) error {
	conn, err := db.Client(c.Context())
	if err != nil {
		logger.Log.Error("Failed to get database connection", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(api_structures.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	c.Locals("db_pool_conn", conn)
	return c.Next()
}
