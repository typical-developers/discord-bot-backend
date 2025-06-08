package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

//	@Router		/guild/{guild_id}/voice-room/lobby/{channel_id}/create [post]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string						true	"The guild ID."
//	@Param		channel_id	path		string						true	"The channel ID."
//	@Param		settings	body		models.VoiceRoomLobbyModify	true	"The activity settings."
//
//	@Success	200			{object}	models.APIResponse[[]models.VoiceRoomLobbyConfig]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	409			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func CreateVoiceRoomLobby(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	var settings *models.VoiceRoomLobbyModify
	if err := c.BodyParser(&settings); err != nil {
		logger.Log.Debug("Failed to parse body.", "error", err)

		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid structure.",
			},
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	// This is a hacky way to set defaults and override them if they're provided.
	// I'm sure there's a clearer way to do this.. will revisit in the future if it has problems.
	creationParams := db.CreateVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: channelId,
		UserLimit:      6,
		CanRename:      false,
		CanLock:        false,
		CanAdjustLimit: false,
	}
	if settings.UserLimit != nil {
		creationParams.UserLimit = *settings.UserLimit
	}
	if settings.CanRename != nil {
		creationParams.CanRename = *settings.CanRename
	}
	if settings.CanLock != nil {
		creationParams.CanLock = *settings.CanLock
	}
	if settings.CanAdjustLimit != nil {
		creationParams.CanAdjustLimit = *settings.CanAdjustLimit
	}

	_, err := queries.CreateVoiceRoomLobby(ctx, creationParams)
	if err != nil {
		errCode, ok := dbutil.UnwrapSQLState(err)
		if ok && errCode == dbutil.SQLStateUniqueViolation {
			logger.Log.Debug("Lobby already exists.", "guild_id", guildId, "error", err)

			return c.Status(fiber.StatusConflict).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "lobby already exists.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to create voice room lobby.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	lobbies, err := queries.GetVoiceRoomLobbies(ctx, guildId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get voice room lobbies.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	mappedLobbies := []models.VoiceRoomLobbyConfig{}
	for _, lobby := range lobbies {
		mappedLobbies = append(mappedLobbies, models.VoiceRoomLobbyConfig{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      int(lobby.UserLimit),
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,
		})
	}

	return c.JSON(models.APIResponse[[]models.VoiceRoomLobbyConfig]{
		Success: true,
		Data:    mappedLobbies,
	})
}

//	@Router		/guild/{guild_id}/voice-room/lobby/{channel_id} [get]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		channel_id	path		string	true	"The channel ID."
//
//	@Success	200			{object}	models.APIResponse[models.VoiceRoomLobbyConfig]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	409			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func GetVoiceRoomLobby(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	lobby, err := queries.GetVoiceRoomLobby(ctx, db.GetVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: channelId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "lobby not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get voice room lobby.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	return c.JSON(models.APIResponse[models.VoiceRoomLobbyConfig]{
		Success: true,
		Data: models.VoiceRoomLobbyConfig{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      int(lobby.UserLimit),
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,
		},
	})
}

//	@Router		/guild/{guild_id}/voice-room/lobby/{channel_id}/update [patch]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string						true	"The guild ID."
//	@Param		channel_id	path		string						true	"The channel ID."
//	@Param		settings	body		models.VoiceRoomLobbyModify	true	"The activity settings."
//
//	@Success	200			{object}	models.APIResponse[models.VoiceRoomLobbyConfig]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	409			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func UpdateVoiceRoomLobby(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	var settings *models.VoiceRoomLobbyModify
	if err := c.BodyParser(&settings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid structure.",
			},
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	_, err := queries.UpdateVoiceRoomLobby(ctx, db.UpdateVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: channelId,
		UserLimit:      dbutil.Int32(settings.UserLimit),
		CanRename:      dbutil.Bool(settings.CanRename),
		CanLock:        dbutil.Bool(settings.CanLock),
		CanAdjustLimit: dbutil.Bool(settings.CanAdjustLimit),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "lobby not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to update voice room lobby.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	lobbies, err := queries.GetVoiceRoomLobbies(ctx, guildId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get voice room lobbies.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}
	mappedLobbies := []models.VoiceRoomLobbyConfig{}
	for _, lobby := range lobbies {
		mappedLobbies = append(mappedLobbies, models.VoiceRoomLobbyConfig{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      int(lobby.UserLimit),
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,
		})
	}

	return c.JSON(models.APIResponse[[]models.VoiceRoomLobbyConfig]{
		Success: true,
		Data:    mappedLobbies,
	})
}

//	@Router		/guild/{guild_id}/voice-room/lobby/{channel_id}/delete [delete]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		channel_id	path		string	true	"The channel ID."
//
//	@Success	200			{object}	models.APIResponse[[]models.VoiceRoomLobbyConfig]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	409			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func DeleteVoiceRoomLobby(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	err := queries.DeleteVoiceRoomLobby(ctx, db.DeleteVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: channelId,
	})
	if err != nil {
		logger.Log.WithSource.Error("Failed to delete voice room lobby.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	lobbies, err := queries.GetVoiceRoomLobbies(ctx, guildId)
	if err != nil {
		logger.Log.WithSource.Error("Failed to get voice room lobbies.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}
	mappedLobbies := []models.VoiceRoomLobbyConfig{}
	for _, lobby := range lobbies {
		mappedLobbies = append(mappedLobbies, models.VoiceRoomLobbyConfig{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      int(lobby.UserLimit),
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,
		})
	}

	return c.JSON(models.APIResponse[[]models.VoiceRoomLobbyConfig]{
		Success: true,
		Data:    mappedLobbies,
	})
}
