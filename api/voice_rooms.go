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

//	@Router		/guild/{guild_id}/voice-room/room/{channel_id}/update [patch]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string					true	"The guild ID."
//	@Param		channel_id	path		string					true	"The channel ID."
//	@Param		settings	body		models.VoiceRoomModify	true	"The activity settings."
//
//	@Success	200			{object}	models.APIResponse[models.VoiceRoomConfig]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	409			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func UpdateVoiceRoom(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	var settings *models.VoiceRoomModify
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

	room, err := queries.UpdateVoiceRoom(ctx, db.UpdateVoiceRoomParams{
		GuildID:        guildId,
		ChannelID:      channelId,
		CurrentOwnerID: dbutil.String(settings.CurrentOwnerID),
		IsLocked:       dbutil.Bool(settings.IsLocked),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "room not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to update voice room.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	return c.JSON(models.APIResponse[models.VoiceRoomConfig]{
		Success: true,
		Data: models.VoiceRoomConfig{
			OriginChannelID: room.OriginChannelID,
			RoomChannelID:   room.ChannelID,
			CreatedByUserID: room.CreatedByUserID,
			CurrentOwnerID:  room.CurrentOwnerID,
			IsLocked:        room.IsLocked.Bool,
		},
	})
}

//	@Router		/guild/{guild_id}/voice-room/room/{channel_id}/unregister [delete]
//	@Tags		Voice Room Lobbies
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path	string	true	"The guild ID."
//	@Param		channel_id	path	string	true	"The channel ID."
//
//	@Success	204
//
//	@Failure	400	{object}	models.APIResponse[ErrorResponse]
//	@Failure	409	{object}	models.APIResponse[ErrorResponse]
//	@Failure	500	{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func DeleteVoiceRoom(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	channelId := c.Params("channel_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	err := queries.DeleteVoiceRoom(ctx, db.DeleteVoiceRoomParams{
		GuildID:   guildId,
		ChannelID: channelId,
	})
	if err != nil {
		logger.Log.WithSource.Error("Failed to delete voice room.", "guild_id", guildId, "channel_id", channelId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
