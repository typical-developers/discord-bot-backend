package handlers

import (
	"github.com/gofiber/fiber/v2"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func SetPoolConn(c *fiber.Ctx) error {
	conn, err := dbutil.Client(c.Context())
	if err != nil {
		logger.Log.WithSource.Error("Failed to get database connection", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	c.Locals("db_pool_conn", conn)
	return c.Next()
}
