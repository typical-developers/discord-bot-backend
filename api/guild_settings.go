package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

//	@Router		/guild/{guild_id}/settings/create [post]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	models.APIResponse[GuildSettings]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func CreateGuildSettings(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := queries.CreateGuildSettings(ctx, guildId)
	if err != nil {
		errCode, ok := dbutil.UnwrapSQLState(err)
		if ok && errCode == dbutil.SQLStateUniqueViolation {
			logger.Log.Debug("Guild already has settings.", "guild_id", guildId, "error", err)

			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild already has settings.",
				},
			})
		}

		logger.Log.Error("Failed to create guild settings.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	return c.JSON(models.APIResponse[models.GuildSettings]{
		Success: true,
		Data: models.GuildSettings{
			ChatActivity: models.ActivityConfig{
				IsEnabled:       settings.ActivityTracking.Bool,
				GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
				CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
				ActivityRoles:   []models.ActivityRole{},
			},
		},
	})
}

//	@Router		/guild/{guild_id}/settings [get]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//
//	@Success	200			{object}	models.APIResponse[GuildSettings]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	404			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func GetGuildSettings(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var mappedChatRoles []models.ActivityRole
	for _, role := range settings.ChatActivityRoles {
		mappedChatRoles = append(mappedChatRoles, models.ActivityRole{
			RoleID:         role.RoleID,
			RequiredPoints: int(role.RequiredPoints.Int32),
		})
	}

	return c.JSON(models.APIResponse[models.GuildSettings]{
		Success: true,
		Data: models.GuildSettings{
			ChatActivity: models.ActivityConfig{
				IsEnabled:       settings.ActivityTracking.Bool,
				GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
				CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
				ActivityRoles:   mappedChatRoles,
			},
		},
	})
}

//	@Router		/guild/{guild_id}/settings/update/activity [patch]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string							true	"The guild ID."
//	@Param		settings	body		models.UpdateActivitySettings	true	"The activity settings."
//
//	@Success	200			{object}	models.APIResponse[GuildSettings]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	404			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func UpdateGuildActivitySettings(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	var activitySettings *models.UpdateActivitySettings
	if err := c.BodyParser(&activitySettings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid structure.",
			},
		})
	}

	println(fmt.Sprintf("activitySettings: %+v", activitySettings))

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	err := queries.UpdateActivitySettings(ctx, db.UpdateActivitySettingsParams{
		GuildID:                  guildId,
		ActivityTracking:         dbutil.Bool(activitySettings.ChatActivity.Enabled),
		ActivityTrackingGrant:    dbutil.Int32(activitySettings.ChatActivity.GrantAmount),
		ActivityTrackingCooldown: dbutil.Int32(activitySettings.ChatActivity.Cooldown),
	})
	if err != nil {
		logger.Log.Error("Failed to update guild settings.", "guild_id", guildId, "error", err, "settings", activitySettings)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		logger.Log.Error("Failed to get guild settings.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	var mappedChatRoles []models.ActivityRole
	for _, role := range settings.ChatActivityRoles {
		mappedChatRoles = append(mappedChatRoles, models.ActivityRole{
			RoleID:         role.RoleID,
			RequiredPoints: int(role.RequiredPoints.Int32),
		})
	}

	return c.JSON(models.APIResponse[models.GuildSettings]{
		Success: true,
		Data: models.GuildSettings{
			ChatActivity: models.ActivityConfig{
				IsEnabled:       settings.ActivityTracking.Bool,
				GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
				CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
				ActivityRoles:   mappedChatRoles,
			},
		},
	})
}

//	@Router		/guild/{guild_id}/settings/update/add-activity-role [patch]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string					true	"The guild ID."
//	@Param		role		body		models.AddActivityRole	true	"The activity settings."
//
//	@Success	200			{object}	models.APIResponse[GuildSettings]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	404			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func GuildAddActivityRole(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	var activityRole models.AddActivityRole
	if err := c.BodyParser(&activityRole); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid structure.",
			},
		})
	}

	if !activityRole.GrantType.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid activity type.",
			},
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	err := queries.InsertActivityRole(ctx, db.InsertActivityRoleParams{
		GuildID:        guildId,
		GrantType:      string(activityRole.GrantType),
		RoleID:         activityRole.RoleID,
		RequiredPoints: int32(activityRole.RequiredPoints),
	})
	if err != nil {
		logger.Log.Error("Failed to add activity role.", "guild_id", guildId, "error", err, "role", activityRole)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		logger.Log.Error("Failed to get guild settings.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	var mappedChatRoles []models.ActivityRole
	for _, role := range settings.ChatActivityRoles {
		mappedChatRoles = append(mappedChatRoles, models.ActivityRole{
			RoleID:         role.RoleID,
			RequiredPoints: int(role.RequiredPoints.Int32),
		})
	}

	return c.JSON(models.APIResponse[models.GuildSettings]{
		Success: true,
		Data: models.GuildSettings{
			ChatActivity: models.ActivityConfig{
				IsEnabled:       settings.ActivityTracking.Bool,
				GrantAmount:     int(settings.ActivityTrackingGrant.Int32),
				CooldownSeconds: int(settings.ActivityTrackingCooldown.Int32),
				ActivityRoles:   mappedChatRoles,
			},
		},
	})
}
