package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/lib/pq"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/internal/pages/layouts"
	u "github.com/typical-developers/discord-bot-backend/internal/usecase"
	discord_state "github.com/typical-developers/discord-bot-backend/pkg/discord-state"
	"maragu.dev/gomponents"
)

type MemberUsecase struct {
	db *sql.DB
	q  *db.Queries
	d  *discord_state.StateManager
}

func NewMemberUsecase(db *sql.DB, q *db.Queries, d *discord_state.StateManager) u.MemberUsecase {
	return &MemberUsecase{db: db, q: q, d: d}
}

func (uc *MemberUsecase) CreateMemberProfile(ctx context.Context, guildId string, userId string) (*u.MemberProfile, error) {
	_, err := uc.d.GuildMember(ctx, guildId, userId)
	if err != nil {
		var dgErr *discordgo.RESTError
		if errors.As(err, &dgErr) && dgErr.Message.Code == discordgo.ErrCodeUnknownMember {
			return nil, u.ErrMemberNotInGuild
		}

		return nil, err
	}

	_, err = uc.q.CreateMemberProfile(ctx, db.CreateMemberProfileParams{
		GuildID:  guildId,
		MemberID: userId,
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, u.ErrMemberProfileExists
		}

		return nil, err
	}

	return uc.GetMemberProfile(ctx, guildId, userId)
}

func (uc *MemberUsecase) GetMemberProfile(ctx context.Context, guildId string, userId string) (*u.MemberProfile, error) {
	guildMember, err := uc.d.GuildMember(ctx, guildId, userId)
	if err != nil {
		var dgErr *discordgo.RESTError
		if errors.As(err, &dgErr) && dgErr.Message.Code == discordgo.ErrCodeUnknownMember {
			return nil, u.ErrMemberNotInGuild
		}

		return nil, err
	}

	chatActivitySettings, err := uc.q.GetGuildChatActivitySettings(ctx, guildId)
	if err != nil {
		return nil, err
	}

	profile, err := uc.q.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: userId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, u.ErrMemberProfileNotFound
		}

		return nil, err
	}

	activityInfo, err := uc.q.GetMemberChatActivityRoleInfo(ctx, db.GetMemberChatActivityRoleInfoParams{
		GuildID: guildId,
		Points:  profile.ChatActivity,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	lastGrant := time.Unix(int64(profile.LastChatActivityGrant), 0)
	nextGrant := lastGrant.Add(time.Duration(chatActivitySettings.GrantCooldown) * time.Second)

	profileInfo := &u.MemberProfile{
		DisplayName: guildMember.DisplayName(),
		Username:    guildMember.User.Username,
		AvatarURL:   guildMember.AvatarURL("100"),

		CardStyle: int32(profile.CardStyle),
		ChatActivity: u.MemberActivity{
			Rank:           int32(profile.ChatActivityRank),
			Points:         int32(profile.ChatActivity),
			LastGrantEpoch: int64(profile.LastChatActivityGrant),
			IsOnCooldown:   time.Now().Before(nextGrant),

			CurrentActivityRoleIds: make([]string, 0),
		},
	}

	if activityInfo.CurrentRoleID.Valid {
		role, err := uc.d.GuildRole(ctx, guildId, activityInfo.CurrentRoleID.String)
		if err == nil {
			profileInfo.ChatActivity.CurrentActivityRole = &u.MemberActivityRole{
				RoleID:         role.ID,
				Name:           role.Name,
				Accent:         fmt.Sprintf("#%06X", role.Color),
				RequiredPoints: int32(activityInfo.CurrentRoleRequiredPoints.Int32),
			}
		}
	}

	if activityInfo.NextRoleID.Valid {
		profileInfo.ChatActivity.NextActivityRole = &u.MemberActivityProgress{
			CurrentProgress:  profile.ChatActivity - activityInfo.CurrentRoleRequiredPoints.Int32,
			RequiredProgress: activityInfo.NextRoleRequiredPoints.Int32 - activityInfo.CurrentRoleRequiredPoints.Int32,
		}
	}

	if len(activityInfo.CurrentRolesIds) > 0 {
		profileInfo.ChatActivity.CurrentActivityRoleIds = activityInfo.CurrentRolesIds
	}

	return profileInfo, nil
}

func (uc *MemberUsecase) IncrementMemberChatActivityPoints(ctx context.Context, guildId string, userId string) (*u.MemberProfile, error) {
	chatActivitySettings, err := uc.q.GetGuildChatActivitySettings(ctx, guildId)
	if err != nil {
		return nil, err
	}

	if !chatActivitySettings.IsEnabled {
		return nil, u.ErrChatActivityTrackingDisabled
	}

	profile, err := uc.q.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: userId,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, u.ErrMemberProfileNotFound
		}

		return nil, err
	}

	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	q := uc.q.WithTx(tx)

	lastGrant := time.Unix(int64(profile.LastChatActivityGrant), 0)
	nextGrant := lastGrant.Add(time.Duration(chatActivitySettings.GrantCooldown) * time.Second)
	if time.Now().Before(nextGrant) {
		return nil, u.ErrMemberOnGrantCooldown
	}

	_, err = q.IncrememberMemberChatActivityPoints(ctx, db.IncrememberMemberChatActivityPointsParams{
		GuildID:  guildId,
		MemberID: userId,

		Points: chatActivitySettings.GrantAmount,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = q.IncrementWeeklyActivityLeaderboard(ctx, db.IncrementWeeklyActivityLeaderboardParams{
		GrantType:    "chat",
		GuildID:      guildId,
		MemberID:     userId,
		EarnedPoints: chatActivitySettings.GrantAmount,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = q.IncrementMonthlyActivityLeaderboard(ctx, db.IncrementMonthlyActivityLeaderboardParams{
		GrantType:    "chat",
		GuildID:      guildId,
		MemberID:     userId,
		EarnedPoints: chatActivitySettings.GrantAmount,
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return uc.GetMemberProfile(ctx, guildId, userId)
}

func (uc *MemberUsecase) GenerateMemberProfileCard(ctx context.Context, guildId string, userId string) (gomponents.Node, error) {
	profile, err := uc.GetMemberProfile(ctx, guildId, userId)
	if err != nil {
		return nil, err
	}

	layout := layouts.ProfileCardProps{
		DisplayName: profile.DisplayName,
		Username:    profile.Username,
		AvatarURL:   profile.AvatarURL,
		ChatActivity: layouts.ActivityInfo{
			Rank:        int(profile.ChatActivity.Rank),
			TotalPoints: int(profile.ChatActivity.Points),
		},
	}

	if profile.ChatActivity.CurrentActivityRole != nil {
		layout.ChatActivity.CurrentTitleInfo = &layouts.ActivityRole{
			Accent: profile.ChatActivity.CurrentActivityRole.Accent,
			Text:   profile.ChatActivity.CurrentActivityRole.Name,
		}
	}

	if profile.ChatActivity.NextActivityRole != nil {
		layout.ChatActivity.RoleCurrentPoints = int(profile.ChatActivity.NextActivityRole.CurrentProgress)
		layout.ChatActivity.RoleRequiredPoints = int(profile.ChatActivity.NextActivityRole.RequiredProgress)
	}

	return layouts.ProfileCard(layout), nil
}

func (uc *MemberUsecase) MigrateMemberProfile(ctx context.Context, guildId string, userId string, toUserId string) error {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := uc.q.WithTx(tx)

	// Fetch the existing profile to migrate from.
	profile, err := q.GetMemberProfile(ctx, db.GetMemberProfileParams{
		GuildID:  guildId,
		MemberID: userId,
	})
	if err != nil {
		_ = tx.Rollback()

		if err == sql.ErrNoRows {
			return u.ErrMemberProfileNotFound
		}

		return err
	}

	// Use the existing profile to mgirate values to the new user.
	// If the member doesn't exist, it will create a new profile with the values.
	// If the member exists, it will update and merge values where necessary (i.e. points).
	err = q.MigrateMemberProfile(ctx, db.MigrateMemberProfileParams{
		GuildID:    guildId,
		ToMemberID: toUserId,
		CardStyle:  profile.CardStyle,

		ChatActivity:          profile.ChatActivity,
		LastChatActivityGrant: profile.LastChatActivityGrant,

		VoiceActivity:          profile.VoiceActivity,
		LastVoiceActivityGrant: profile.LastVoiceActivityGrant,
	})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Reset the old profile to default values.
	err = q.ResetMemberProfile(ctx, db.ResetMemberProfileParams{
		GuildID:  guildId,
		MemberID: userId,
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
