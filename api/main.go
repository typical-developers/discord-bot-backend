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
		guildSnowflakeHandler := handlers.CheckSnowflakeParams([]string{"guild_id"})

		guild.Get("/activity-leaderboard/card", guildSnowflakeHandler, ActivityLeaderboardCard)

		guild.Post("/create-settings", guildSnowflakeHandler, CreateGuildSettings)
		guild.Get("/settings", guildSnowflakeHandler, GetGuildSettings)
		guild.Patch("/update-settings", guildSnowflakeHandler, UpdateGuildSettings)

		member := guild.Group("/member/:member_id")
		{
			memberSnowflakeHandler := handlers.CheckSnowflakeParams([]string{"guild_id", "member_id"})

			member.Post("/create-profile", memberSnowflakeHandler, CreateMemberProfile)
			member.Get("/profile", memberSnowflakeHandler, GetMemberProfile)
			member.Get("/profile/card", memberSnowflakeHandler, MemberProfileCard)
			member.Post("/profile/increment-points", memberSnowflakeHandler, IncrementActivityPoints)
		}
	}

	// All HTML related assets are publicly accessible if the endpoint is known.
	// They're not the most sensitive thing and securing them would be annoying.
	app.Get("/html/*", GetHTMLAsset)
	app.Get("/html-version.json", Version)
}
