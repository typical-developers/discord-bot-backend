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
//	@tag.name					Voice Room Lobbies
//	@tag.description			Voice room lobby endpoints.
//
//	@tag.name					Voice Rooms
//	@tag.description			Voice room endpoints.
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

		guild.Get("/settings", guildSnowflakeHandler, GetGuildSettings)
		guild.Post("/settings/create", guildSnowflakeHandler, CreateGuildSettings)
		guild.Patch("/settings/update/activity", guildSnowflakeHandler, UpdateGuildActivitySettings)
		guild.Post("/settings/update/add-activity-role", guildSnowflakeHandler, GuildAddActivityRole)
		guild.Patch("/member-migrate", guildSnowflakeHandler, MigrateMemberProfile)

		memberSnowflakeHandler := handlers.CheckSnowflakeParams([]string{"guild_id", "member_id"})
		member := guild.Group("/member/:member_id", memberSnowflakeHandler)
		{
			member.Get("/profile", GetMemberProfile)
			member.Post("/profile/create", CreateMemberProfile)
			member.Get("/profile/card", MemberProfileCard)
			member.Post("/profile/increment-points", IncrementActivityPoints)
		}

		voiceRoomSnowflakeHandler := handlers.CheckSnowflakeParams([]string{"guild_id", "channel_id"})
		voiceRoom := guild.Group("/voice-room")
		{
			lobby := voiceRoom.Group("/lobby/:channel_id", voiceRoomSnowflakeHandler)
			{
				lobby.Post("/create", CreateVoiceRoomLobby)
				lobby.Post("/register", RegisterVoiceRoom)
				lobby.Get("/", GetVoiceRoomLobby)
				lobby.Patch("/update", UpdateVoiceRoomLobby)
				lobby.Delete("/delete", DeleteVoiceRoomLobby)
			}

			room := voiceRoom.Group("/room/:channel_id", voiceRoomSnowflakeHandler)
			{
				room.Delete("/unregister", DeleteVoiceRoom)
				room.Patch("/update", UpdateVoiceRoom)
			}
		}
	}

	// All HTML related assets are publicly accessible if the endpoint is known.
	// They're not the most sensitive thing and securing them would be annoying.
	app.Get("/html/*", GetHTMLAsset)
	app.Get("/html-version.json", Version)
}
