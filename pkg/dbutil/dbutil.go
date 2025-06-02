package dbutil

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	models "github.com/typical-developers/discord-bot-backend/internal"
	"github.com/typical-developers/discord-bot-backend/internal/db"
)

// Errors to help identify SQLState codes more cleanly.
type SQLState string

const (
	SQLStateUniqueViolation SQLState = "23505"
)

// Utility to unwrap a pgconn.PgError cleanly.
func UnwrapSQLState(err error) (SQLState, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return SQLState(pgErr.Code), true
	}

	return "", false
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
	Next    models.ActivityRoleProgress
	Current []models.ActivityRole
}

// Gets information on the member's next role and current roles based on their points.
// Roles must be fetched and provided separately.
func MapMemberRoles(points int, activityRoles []db.GetGuildActivityRolesRow) MemberRoles {
	nextRole := models.ActivityRoleProgress{}
	currentRoles := []models.ActivityRole{}

	requiredPoints := int32(0)
	for _, role := range activityRoles {
		requiredPoints += role.RequiredPoints.Int32
		if requiredPoints <= int32(points) {
			currentRoles = append(currentRoles, models.ActivityRole{
				RoleID:         role.RoleID,
				RequiredPoints: int(role.RequiredPoints.Int32),
			})

			continue
		}

		nextRole = models.ActivityRoleProgress{
			RoleID:          role.RoleID,
			Progress:        points - (int(requiredPoints - role.RequiredPoints.Int32)),
			RemainingPoints: int(requiredPoints) - points,
			RequiredPoints:  int(role.RequiredPoints.Int32),
		}
		break
	}

	return MemberRoles{
		Next:    nextRole,
		Current: currentRoles,
	}
}
