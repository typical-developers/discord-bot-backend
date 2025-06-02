package api

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
	"github.com/typical-developers/discord-bot-backend/pkg/regexutil"
)

//	@Router		/guild/{guild_id}/member/{member_id}/create [post]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		member_id	path		string	true	"The member ID."
//
//	@Success	200			{object}	models.MemberProfile
//
//	@Failure	400			{object}	models.GenericResponse
//	@Failure	500			{object}	models.GenericResponse
//
// nolint:staticcheck
func CreateMemberProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")

	if !regexutil.Snowflake.MatchString(guildId) {
		return c.Status(fiber.StatusBadRequest).JSON(models.GenericResponse{
			Success: false,
			Message: "guild_id is not snowflake.",
		})
	}

	if !regexutil.Snowflake.MatchString(memberId) {
		return c.Status(fiber.StatusBadRequest).JSON(models.GenericResponse{
			Success: false,
			Message: "member_id is not snowflake.",
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.Error("Failed to create member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}
	queries := db.New(connection).WithTx(tx)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.GenericResponse{
				Success: false,
				Message: "guild settings not found.",
			})
		}

		logger.Log.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := queries.CreateMemberProfile(ctx, db.CreateMemberProfileParams{
		GuildID:        guildId,
		MemberID:       memberId,
		ActivityPoints: 0,
	})
	if err != nil {
		logger.Log.Error("Failed to create member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		_ = tx.Rollback(ctx)

		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		logger.Log.Error("Failed to get member rankings.", "guild_id", guildId, "member_id", memberId, "error", err)
		_ = tx.Rollback(ctx)

		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.Error("Failed to commit transaction.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)

	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	return c.JSON(models.MemberProfile{
		CardStyle: models.CardStyle(profile.CardStyle),
		ChatActivity: models.MemberActivity{
			Rank:         int(rankings.ChatRank),
			LastGrant:    grantTime,
			IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, int(settings.ActivityTrackingCooldown.Int32)),
			Points:       int(profile.ActivityPoints),
			Roles: models.MemberRoles{
				Next:     roles.Next,
				Obtained: roles.Current,
			},
		},
	})
}

//	@Router		/guild/{guild_id}/member/{member_id} [get]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		member_id	path		string	true	"The member ID."
//
//	@Success	200			{object}	models.MemberProfile
//
//	@Failure	400			{object}	models.GenericResponse
//	@Failure	500			{object}	models.GenericResponse
//
// nolint:staticcheck
func GetMemberProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.GenericResponse{
				Success: false,
				Message: "guild settings not found.",
			})
		}

		logger.Log.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := dbutil.GetMemberProfile(ctx, queries, guildId, memberId)
	if err != nil {
		logger.Log.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)

		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "failed to get member profile.",
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)

	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	return c.JSON(models.MemberProfile{
		CardStyle: models.CardStyle(profile.CardStyle),
		ChatActivity: models.MemberActivity{
			Rank:         int(profile.ChatRank),
			LastGrant:    grantTime,
			IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, int(settings.ActivityTrackingCooldown.Int32)),
			Points:       int(profile.ActivityPoints),
			Roles: models.MemberRoles{
				Next:     roles.Next,
				Obtained: roles.Current,
			},
		},
	})
}

//	@Router		/guild/{guild_id}/member/{member_id}/activity-points/increment [post]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path		string				true	"The guild ID."
//	@Param		member_id		path		string				true	"The member ID."
//	@Param		activity_type	query		models.ActivityType	true	"The activity type."
//
//	@Success	200				{object}	models.MemberProfile
//
//	@Failure	400				{object}	models.GenericResponse
//	@Failure	500				{object}	models.GenericResponse
//
// nolint:staticcheck
func IncrementActivityPoints(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")
	activityType := models.ActivityType(c.Query("activity_type"))

	if !regexutil.Snowflake.MatchString(guildId) {
		return c.Status(fiber.StatusBadRequest).JSON(models.GenericResponse{
			Success: false,
			Message: "guild_id is not snowflake.",
		})
	}

	if !regexutil.Snowflake.MatchString(memberId) {
		return c.Status(fiber.StatusBadRequest).JSON(models.GenericResponse{
			Success: false,
			Message: "member_id is not snowflake.",
		})
	}

	if !activityType.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(models.GenericResponse{
			Success: false,
			Message: "activity_type is not valid.",
		})
	}

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	queries := db.New(connection)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.GenericResponse{
				Success: false,
				Message: "guild settings not found.",
			})
		}

		logger.Log.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := queries.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		logger.Log.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	lastGrant := time.Unix(int64(profile.LastGrantEpoch), 0)
	cooldown := int(settings.ActivityTrackingCooldown.Int32)
	if activityType == models.ActivityTypeChat {
		if dbutil.IsMemberOnCooldown(lastGrant, cooldown) {
			return c.Status(fiber.StatusForbidden).JSON(models.GenericResponse{
				Success: false,
				Message: "cannot grant, member is on cooldown.",
			})
		}

		updatedProfile, err := queries.IncrememberMemberChatActivityPoints(ctx, db.IncrememberMemberChatActivityPointsParams{
			Points:   settings.ActivityTrackingGrant.Int32,
			GuildID:  guildId,
			MemberID: memberId,
		})
		if err != nil {
			logger.Log.Error("Failed to increment member chat activity points.", "guild_id", guildId, "member_id", memberId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
				Success: false,
				Message: "internal server error.",
			})
		}

		profile = db.GetMemberProfileRow{
			CardStyle:      updatedProfile.CardStyle,
			ActivityPoints: updatedProfile.ActivityPoints,
			LastGrantEpoch: updatedProfile.LastGrantEpoch,
		}
	}

	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		logger.Log.Error("Failed to get member rankings.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.GenericResponse{
			Success: false,
			Message: "internal server error.",
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)

	return c.JSON(models.MemberProfile{
		CardStyle: models.CardStyle(profile.CardStyle),
		ChatActivity: models.MemberActivity{
			Rank:         int(rankings.ChatRank),
			LastGrant:    grantTime,
			IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, cooldown),
			Points:       int(profile.ActivityPoints),
			Roles: models.MemberRoles{
				Next:     roles.Next,
				Obtained: roles.Current,
			},
		},
	})
}
