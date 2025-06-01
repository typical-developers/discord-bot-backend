package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/typical-developers/discord-bot-backend/api"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	_ "github.com/typical-developers/discord-bot-backend/internal/docs"
)

func main() {
	_, err := db.InitalizePool()
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(logger.New())

	api.Register(app)

	_ = app.Listen(":8080")
}
