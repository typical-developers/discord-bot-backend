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
//	@tag.name					Members
//	@tag.description			Member endpoints.
//
//	@tag.name					HTML Generation
//	@tag.description			HTML generation endpoints.
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

	guild := app.Group("/guild/:guild_id", handlers.CheckAuthorization, handlers.LogRequest, handlers.SetPoolConn)
	{
		guild.Post("/create", CreateGuildSettings)
		guild.Get("/", GetGuildSettings)
		guild.Patch("/update", UpdateGuildSettings)

		member := guild.Group("/member/:member_id")
		{
			member.Post("/create", CreateMemberProfile)
			member.Get("/", GetMemberProfile)
			member.Get("/profile-card", MemberProfileCard)
			member.Post("/activity-points/increment", IncrementActivityPoints)
		}
	}

	// All HTML related assets are publicly accessible if the endpoint is known.
	// They're not the most sensitive thing and securing them would be annoying.
	app.Get("/html/*", GetHTMLAsset)
	app.Get("/html-version.json", Version)
}
