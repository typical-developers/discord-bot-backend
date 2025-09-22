package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/typical-developers/discord-bot-backend/internal/db"

	u "github.com/typical-developers/discord-bot-backend/internal/usecase"
)

type GuildUsecase struct {
	q db.Querier
}

func NewGuildUsecase(q db.Querier) u.GuildsUsecase {
	return &GuildUsecase{q: q}
}

func (uc *GuildUsecase) CreateGuildSettings(ctx context.Context, guildId string) (*u.GuildSettings, error) {
	_, err := uc.q.CreateGuildSettings(ctx, guildId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr); pqErr.Code == "23505" {
			return nil, u.ErrGuildSettingsExist
		}

		return nil, err
	}

	return uc.GetGuildSettings(ctx, guildId)
}

func (uc *GuildUsecase) GetGuildSettings(ctx context.Context, guildId string) (*u.GuildSettings, error) {
	settings, err := uc.q.GetGuildSettings(ctx, guildId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, u.ErrGuildNotFound
		}

		return nil, err
	}

	chatActivityRoles, err := uc.q.GetGuildActivityRoles(ctx, db.GetGuildActivityRolesParams{
		GuildID:      guildId,
		ActivityType: "chat",
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	chatRoles := make([]u.GuildActivityRole, 0)
	for _, role := range chatActivityRoles {
		chatRoles = append(chatRoles, u.GuildActivityRole{
			RoleID:         role.RoleID,
			RequiredPoints: role.RequiredPoints.Int32,
		})
	}

	creationLobbies, err := uc.q.GetVoiceRoomLobbies(ctx, guildId)
	if err != nil {
		return nil, err
	}

	lobbies := make([]u.VoiceRoomLobby, 0)
	for _, lobby := range creationLobbies {
		lobbies = append(lobbies, u.VoiceRoomLobby{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      lobby.UserLimit,
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,
		})
	}

	return &u.GuildSettings{
		ChatActivityTracking: u.GuildActivityTracking{
			Enabled:    settings.ChatActivityTracking.Valid,
			Cooldown:   settings.ChatActivityGrant.Int32,
			Grant:      settings.ChatActivityGrant.Int32,
			GrantRoles: chatRoles,
			DenyRoles:  []string{},
		},
		VoiceRoomLobbies: lobbies,
	}, nil
}
