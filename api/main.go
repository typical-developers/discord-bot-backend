package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/typical-developers/discord-bot-backend/handlers"
)

//	@title						Discord Bot API
//	@version					1.0
//	@description				The API for the main Typical Developers Discord bot.
//
//	@tag.name					Guilds
//	@tag.description			Guild endpoints.
//
//	@securitydefinitions.apikey	APIKeyAuth
//	@in							header
//	@name						X-API-KEY
//
// nolint:staticcheck
func Register(app *fiber.App) {
	// Registers swagger documentation.
	app.Get("/docs/*", swagger.New(swagger.Config{
		DeepLinking:  false,
		DocExpansion: "list",
	}))

	guildSettings := app.Group("/guild-settings", handlers.CheckAuthorization, handlers.LogRequest, handlers.SetPoolConn)
	{
		guildSettings.Post("/:guild_id/create",
			CreateGuildSettings,
		)

		guildSettings.Get("/:guild_id",
			GetGuildSettings,
		)

		guildSettings.Patch("/:guild_id/update",
			UpdateGuildSettings,
		)
	}
}
