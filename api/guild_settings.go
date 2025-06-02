package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	api_structures "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
	"github.com/typical-developers/discord-bot-backend/pkg/regexutil"
)

//	@Router		/guild/{guild_id}/create [post]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	api_structures.GuildSettings
//
//	@Failure	400			{object}	api_structures.GenericResponse
//	@Failure	500			{object}	api_structures.GenericResponse
//
// nolint:staticcheck
func CreateGuildSettings(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	if !regexutil.Snowflake.MatchString(guildId) {
		return c.Status(fiber.StatusBadRequest).JSON(api_structures.GenericResponse{
			Success: false,
			Message: "guild_id is not snowflake.",
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := queries.CreateGuildSettings(ctx, guildId)
	if err != nil {
		errCode, ok := dbutil.UnwrapSQLState(err)
		if ok && errCode == dbutil.SQLStateUniqueViolation {
			logger.Log.Debug("Guild already has settings.", "guild_id", guildId, "error", err)

			return c.Status(fiber.StatusBadRequest).JSON(api_structures.GenericResponse{
				Success: false,
				Message: "guild already has settings.",
			})
		}

		logger.Log.Error("Failed to create guild settings.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(api_structures.GenericResponse{
			Success: false,
			Message: "failed to create guild settings.",
		})
	}

	return c.JSON(api_structures.GuildSettings{
		ChatActivity: api_structures.ActivityConfig{
			IsEnabled:       settings.ActivityTracking.Bool,
			GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
			CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
			ActivityRoles:   []api_structures.ActivityRole{},
		},
	})
}

//	@Router		/guild/{guild_id} [get]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	api_structures.GuildSettings
//
//	@Failure	400			{object}	api_structures.GenericResponse
//	@Failure	404			{object}	api_structures.GenericResponse
//	@Failure	500			{object}	api_structures.GenericResponse
//
// nolint:staticcheck
func GetGuildSettings(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	if !regexutil.Snowflake.MatchString(guildId) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(api_structures.GenericResponse{
				Success: false,
				Message: "guild settings not found.",
			})
		}

		logger.Log.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var mappedChatRoles []api_structures.ActivityRole
	for _, role := range settings.ChatActivityRoles {
		mappedChatRoles = append(mappedChatRoles, api_structures.ActivityRole{
			RoleID:         role.RoleID,
			RequiredPoints: int(role.RequiredPoints.Int32),
		})
	}

	return c.JSON(api_structures.GuildSettings{
		ChatActivity: api_structures.ActivityConfig{
			IsEnabled:       settings.ActivityTracking.Bool,
			GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
			CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
			ActivityRoles:   mappedChatRoles,
		},
	})
}

//	@Router		/guild/{guild_id}/update [patch]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	api_structures.GuildSettings
//
//	@Failure	400			{object}	api_structures.GenericResponse
//	@Failure	404			{object}	api_structures.GenericResponse
//	@Failure	500			{object}	api_structures.GenericResponse
//
// nolint:staticcheck
func UpdateGuildSettings(c *fiber.Ctx) error {
	return nil
}
