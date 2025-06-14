package dbutil

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
)

// Get a client from the pool.
func Client(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := db.Pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// ---------------------------------------------------------------------

// Adds cooldownSeconds to lastGrant to determine if the member is on a cooldown.
func IsMemberOnCooldown(lastGrant time.Time, cooldownSeconds int) bool {
	nextGrant := lastGrant.Add(time.Duration(cooldownSeconds) * time.Second)
	return nextGrant.After(time.Now())
}

type GuildSettings struct {
	db.GetGuildSettingsRow
	ChatActivityRoles []db.GetGuildActivityRolesRow
}

// Fetches all relating guild settings.
// This includes configurations and activity roles.
func GetGuildSettings(ctx context.Context, queries *db.Queries, guildId string) (*GuildSettings, error) {
	settings, err := queries.GetGuildSettings(ctx, guildId)
	if err != nil {
		return nil, err
	}

	chatGrantRoles, err := queries.GetGuildActivityRoles(ctx, db.GetGuildActivityRolesParams{
		GuildID:      guildId,
		ActivityType: string(models.ActivityTypeChat),
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	return &GuildSettings{
		GetGuildSettingsRow: settings,
		ChatActivityRoles:   chatGrantRoles,
	}, nil
}

type MemberProfile struct {
	db.GetMemberProfileRow
	db.GetMemberRankingsRow
}

// Fetches all relating member profile information.
func GetMemberProfile(ctx context.Context, queries *db.Queries, guildId string, memberId string) (*MemberProfile, error) {
	profile, err := queries.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		return nil, err
	}

	rankings, err := queries.GetMemberRankings(ctx, db.GetMemberRankingsParams{
		GuildID:  guildId,
		MemberID: memberId,
	})
	if err != nil {
		return nil, err
	}

	return &MemberProfile{
		GetMemberProfileRow:  profile,
		GetMemberRankingsRow: rankings,
	}, nil
}

type MemberRoles struct {
	Next     *models.ActivityRoleProgress
	Current  *models.ActivityRole
	Obtained []models.ActivityRole
}

// Gets information on the member's next role and current roles based on their points.
// Roles must be fetched and provided separately.
func MapMemberRoles(points int, activityRoles []db.GetGuildActivityRolesRow) MemberRoles {
	obtainedRoles := []models.ActivityRole{}
	var nextRole *models.ActivityRoleProgress
	var currentRole *models.ActivityRole

	lastRequired := int32(0)
	for _, role := range activityRoles {
		rolePoints := role.RequiredPoints.Int32

		if rolePoints <= int32(points) {
			obtainedRole := models.ActivityRole{
				RoleID:         role.RoleID,
				RequiredPoints: int(rolePoints),
			}

			currentRole = &obtainedRole
			lastRequired = rolePoints

			obtainedRoles = append(obtainedRoles, obtainedRole)

			continue
		}

		nextRole = &models.ActivityRoleProgress{
			RoleID:         role.RoleID,
			Progress:       points - int(lastRequired),
			RequiredPoints: int(rolePoints - lastRequired),
		}

		break
	}

	return MemberRoles{
		Next:     nextRole,
		Current:  currentRole,
		Obtained: obtainedRoles,
	}
}
