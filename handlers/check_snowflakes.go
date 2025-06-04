package handlers

import (
	"github.com/gofiber/fiber/v2"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/pkg/regexutil"
)

func CheckSnowflakeParams(key []string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		snowflakes := []regexutil.SnowflakeIDs{}

		for _, k := range key {
			snowflakes = append(snowflakes, regexutil.SnowflakeIDs{
				Key: k,
				ID:  c.Params(k),
			})
		}

		if err := regexutil.CheckSnowflakes(snowflakes); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: err.Error(),
				},
			})
		}

		return c.Next()
	}
}
