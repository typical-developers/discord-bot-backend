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

//	@Router		/guild/{guild_id}/member/{member_id}/profile/create [post]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		member_id	path		string	true	"The member ID."
//
//	@Success	200			{object}	models.APIResponse[MemberProfile]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func CreateMemberProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.WithSource.Error("Failed to create member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}
	queries := db.New(connection).WithTx(tx)
	defer connection.Release()

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := queries.CreateMemberProfile(ctx, db.CreateMemberProfileParams{
		GuildID:        guildId,
		MemberID:       memberId,
		ActivityPoints: 0,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		errCode, ok := dbutil.UnwrapSQLState(err)
		if ok && errCode == dbutil.SQLStateUniqueViolation {
			logger.Log.Debug("Guild already has settings.", "guild_id", guildId, "error", err)

			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "member already has profile.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to create member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		logger.Log.WithSource.Error("Failed to get member rankings.", "guild_id", guildId, "member_id", memberId, "error", err)
		_ = tx.Rollback(ctx)

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.WithSource.Error("Failed to commit transaction.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)

	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	return c.JSON(models.APIResponse[models.MemberProfile]{
		Success: true,
		Data: models.MemberProfile{
			CardStyle: models.CardStyle(profile.CardStyle),
			ChatActivity: models.MemberActivity{
				Rank:         int(rankings.ChatRank),
				LastGrant:    grantTime,
				IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, int(settings.ActivityTrackingCooldown.Int32)),
				Points:       int(profile.ActivityPoints),
				Roles: models.MemberRoles{
					Next: roles.Next,
					// Obtained: roles.Current,
				},
			},
		},
	})
}

//	@Router		/guild/{guild_id}/member/{member_id}/profile [get]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string	true	"The guild ID."
//	@Param		member_id	path		string	true	"The member ID."
//
//	@Success	200			{object}	models.APIResponse[MemberProfile]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
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
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := dbutil.GetMemberProfile(ctx, queries, guildId, memberId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "member not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)

		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)
	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	return c.JSON(models.APIResponse[models.MemberProfile]{
		Success: true,
		Data: models.MemberProfile{
			CardStyle: models.CardStyle(profile.CardStyle),
			ChatActivity: models.MemberActivity{
				Rank:         int(profile.ChatRank),
				LastGrant:    grantTime,
				IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, int(settings.ActivityTrackingCooldown.Int32)),
				Points:       int(profile.ActivityPoints),
				Roles: models.MemberRoles{
					Next:     roles.Next,
					Obtained: roles.Obtained,
				},
			},
		},
	})
}

//	@Router		/guild/{guild_id}/member/{member_id}/profile/increment-points [post]
//	@Tags		Members
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id		path		string				true	"The guild ID."
//	@Param		member_id		path		string				true	"The member ID."
//	@Param		activity_type	query		models.ActivityType	true	"The activity type."
//
//	@Success	200				{object}	models.APIResponse[MemberActivity]
//
//	@Failure	400				{object}	models.APIResponse[ErrorResponse]
//	@Failure	500				{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func IncrementActivityPoints(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")
	memberId := c.Params("member_id")
	activityType := models.ActivityType(c.Query("activity_type"))

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	defer connection.Release()

	if !activityType.Valid() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "activity_type is not valid (chat).",
			},
		})
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.WithSource.Error("Failed to start transaction.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}
	queries := db.New(connection).WithTx(tx)

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	profile, err := queries.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "member not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get member profile.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	lastGrant := time.Unix(int64(profile.LastGrantEpoch), 0)
	cooldown := int(settings.ActivityTrackingCooldown.Int32)
	if activityType == models.ActivityTypeChat {
		if !settings.ActivityTracking.Bool {
			_ = tx.Rollback(ctx)

			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "activity tracking is not enabled.",
				},
			})
		}

		if dbutil.IsMemberOnCooldown(lastGrant, cooldown) {
			_ = tx.Rollback(ctx)

			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "cannot grant, member is on cooldown.",
				},
			})
		}

		updatedProfile, err := queries.IncrememberMemberChatActivityPoints(ctx, db.IncrememberMemberChatActivityPointsParams{
			Points:   settings.ActivityTrackingGrant.Int32,
			GuildID:  guildId,
			MemberID: memberId,
		})
		if err != nil {
			_ = tx.Rollback(ctx)

			logger.Log.WithSource.Error("Failed to increment member chat activity points.", "guild_id", guildId, "member_id", memberId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}

		profile = db.GetMemberProfileRow{
			CardStyle:      updatedProfile.CardStyle,
			ActivityPoints: updatedProfile.ActivityPoints,
			LastGrantEpoch: updatedProfile.LastGrantEpoch,
		}

		err = queries.IncrementWeeklyActivityLeaderboard(ctx, db.IncrementWeeklyActivityLeaderboardParams{
			GrantType:    string(models.ActivityTypeChat),
			GuildID:      guildId,
			MemberID:     memberId,
			EarnedPoints: int32(settings.ActivityTrackingGrant.Int32),
		})
		if err != nil {
			_ = tx.Rollback(ctx)

			logger.Log.WithSource.Error("Failed to increment weekly activity leaderboard.", "guild_id", guildId, "member_id", memberId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}

		err = queries.IncrementMonthlyActivityLeaderboard(ctx, db.IncrementMonthlyActivityLeaderboardParams{
			GrantType:    string(models.ActivityTypeChat),
			GuildID:      guildId,
			MemberID:     memberId,
			EarnedPoints: int32(settings.ActivityTrackingGrant.Int32),
		})
		if err != nil {
			_ = tx.Rollback(ctx)

			logger.Log.WithSource.Error("Failed to increment monthly activity leaderboard.", "guild_id", guildId, "member_id", memberId, "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "internal server error.",
				},
			})
		}
	}

	grantTime := time.Unix(int64(profile.LastGrantEpoch), 0)
	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		logger.Log.WithSource.Error("Failed to get member rankings.", "guild_id", guildId, "member_id", memberId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	_ = tx.Commit(ctx)

	roles := dbutil.MapMemberRoles(int(profile.ActivityPoints), settings.ChatActivityRoles)
	return c.JSON(models.APIResponse[models.MemberActivity]{
		Success: true,
		Data: models.MemberActivity{
			Rank:         int(rankings.ChatRank),
			LastGrant:    grantTime,
			IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, cooldown),
			Points:       int(profile.ActivityPoints),
			Roles: models.MemberRoles{
				Next:     roles.Next,
				Obtained: roles.Obtained,
			},
		},
	})
}

//	@Router		/guild/{guild_id}/member-migrate [patch]
//	@Tags		Guilds
//
//	@Security	APIKeyAuth
//
//	@Param		guild_id	path		string					true	"The guild ID."
//	@Param		ids			body		models.MigrateProfile	true	"The members to transfer points from -> to."
//
//	@Success	200			{object}	models.APIResponse[MemberActivity]
//
//	@Failure	400			{object}	models.APIResponse[ErrorResponse]
//	@Failure	500			{object}	models.APIResponse[ErrorResponse]
//
// nolint:staticcheck
func MigrateMemberProfile(c *fiber.Ctx) error {
	ctx := c.Context()
	guildId := c.Params("guild_id")

	connection := c.Locals("db_pool_conn").(*pgxpool.Conn)
	defer connection.Release()

	var info *models.MigrateProfile
	if err := c.BodyParser(&info); err != nil {
		logger.Log.Debug("Failed to parse body.", "error", err)

		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "invalid structure.",
			},
		})
	}

	// We only check if to_id is a proper snowflake.
	// Since we check if the from_id exists a bit down, we don't need to worry about it being a snowflake or not.
	// to_id actually has a chance to *create* a resource.
	if err := regexutil.CheckSnowflake(info.ToID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "to_id is not a valid snowflake.",
			},
		})
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.WithSource.Error("Failed to start transaction.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}
	queries := db.New(connection)
	queriesTx := queries.WithTx(tx)

	settings, err := dbutil.GetGuildSettings(ctx, queries, guildId)
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "guild settings not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get guild settings", "error", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// --- Migration Starts

	// First, we make sure the fromId exists.
	_, err = queriesTx.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: info.FromID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(models.APIResponse[models.ErrorResponse]{
				Success: false,
				Data: models.ErrorResponse{
					Message: "from member not found.",
				},
			})
		}

		logger.Log.WithSource.Error("Failed to get member profile.", "guild_id", guildId, "member_id", info.FromID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	// Then, we migrate from_id profile -> to_id profile.
	newProfile, err := queriesTx.MigrateMemberProfile(ctx, db.MigrateMemberProfileParams{
		GuildID: guildId,
		FromID:  info.FromID,
		ToID:    info.ToID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		logger.Log.WithSource.Error("Failed to migrate member profile.", "guild_id", guildId, "from_id", info.FromID, "to_id", info.ToID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	// Finally, we reset the old profile to defaults.
	err = queriesTx.ResetOldMemberProfile(ctx, db.ResetOldMemberProfileParams{
		GuildID:  guildId,
		MemberID: info.FromID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		logger.Log.WithSource.Error("Failed to reset old member profile.", "guild_id", guildId, "member_id", info.FromID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.WithSource.Error("Failed to commit transaction.", "guild_id", guildId, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	// --- Migration Ends

	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: info.ToID,
	})
	if err != nil {
		_ = tx.Rollback(ctx)

		logger.Log.WithSource.Error("Failed to get member profile rankings.", "guild_id", guildId, "member_id", info.ToID, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse[models.ErrorResponse]{
			Success: false,
			Data: models.ErrorResponse{
				Message: "internal server error.",
			},
		})
	}

	roles := dbutil.MapMemberRoles(int(newProfile.ActivityPoints), settings.ChatActivityRoles)
	grantTime := time.Unix(int64(newProfile.LastGrantEpoch), 0)
	return c.JSON(models.APIResponse[models.MemberProfile]{
		Success: true,
		Data: models.MemberProfile{
			CardStyle: models.CardStyle(newProfile.CardStyle),
			ChatActivity: models.MemberActivity{
				Rank:         int(rankings.ChatRank),
				LastGrant:    grantTime,
				IsOnCooldown: dbutil.IsMemberOnCooldown(grantTime, int(settings.ActivityTrackingCooldown.Int32)),
				Points:       int(newProfile.ActivityPoints),
				Roles: models.MemberRoles{
					Next:     roles.Next,
					Obtained: roles.Obtained,
				},
			},
		},
	})
}
