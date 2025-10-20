package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/internal/pages/layouts"
	discord_state "github.com/typical-developers/discord-bot-backend/pkg/discord-state"
	"github.com/typical-developers/discord-bot-backend/pkg/sqlx"
	"maragu.dev/gomponents"

	u "github.com/typical-developers/discord-bot-backend/internal/usecase"
)

type GuildUsecase struct {
	db *sql.DB
	q  *db.Queries
	d  *discord_state.StateManager
}

func NewGuildUsecase(db *sql.DB, q *db.Queries, d *discord_state.StateManager) u.GuildsUsecase {
	return &GuildUsecase{db: db, q: q, d: d}
}

func (uc *GuildUsecase) CreateGuildSettings(ctx context.Context, guildId string) (*u.GuildSettings, error) {
	_, err := uc.q.CreateGuildSettings(ctx, guildId)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr); pqErr.Code == "23505" {
			return nil, u.ErrGuildSettingsExists
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, u.ErrGuildNotFound
		}

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
		currentRooms, err := uc.q.GetVoiceRooms(ctx, db.GetVoiceRoomsParams{
			GuildID:         guildId,
			OriginChannelID: lobby.VoiceChannelID,
		})

		if err != nil {
			return nil, err
		}

		rooms := make([]u.VoiceRoom, 0)
		for _, room := range currentRooms {
			rooms = append(rooms, u.VoiceRoom{
				OriginChannelId: room.OriginChannelID,
				CreatorId:       room.CreatedByUserID,
				CurrentOwnerId:  room.CurrentOwnerID,
				IsLocked:        room.IsLocked.Valid && room.IsLocked.Bool,
			})
		}

		lobbies = append(lobbies, u.VoiceRoomLobby{
			ChannelID:      lobby.VoiceChannelID,
			UserLimit:      lobby.UserLimit,
			CanRename:      lobby.CanRename,
			CanLock:        lobby.CanLock,
			CanAdjustLimit: lobby.CanAdjustLimit,

			VoiceRooms: rooms,
		})
	}

	return &u.GuildSettings{
		ChatActivityTracking: u.GuildActivityTracking{
			IsEnabled:       settings.ChatActivityTracking.Bool,
			CooldownSeconds: settings.ChatActivityCooldown.Int32,
			GrantAmount:     settings.ChatActivityGrant.Int32,
			ActivityRoles:   chatRoles,
			DenyRoles:       []string{},
		},
		VoiceRoomLobbies: lobbies,
	}, nil
}

func (uc *GuildUsecase) UpdateGuildActivitySettings(ctx context.Context, guildId string, opts u.UpdateAcitivtySettings) (*u.GuildSettings, error) {
	err := uc.q.UpdateActivitySettings(ctx, db.UpdateActivitySettingsParams{
		GuildID:              guildId,
		ChatActivityTracking: sqlx.Bool(opts.ChatActivity.IsEnabled),
		ChatActivityGrant:    sqlx.Int32(opts.ChatActivity.GrantAmount),
		ChatActivityCooldown: sqlx.Int32(opts.ChatActivity.CooldownSeconds),
	})
	if err != nil {
		return nil, err
	}

	return uc.GetGuildSettings(ctx, guildId)
}

func (uc *GuildUsecase) CreateActivityRole(ctx context.Context, guildId string, activityType string, roleId string, requiredPoints int32) (*u.GuildActivityRole, error) {
	err := uc.q.InsertActivityRole(ctx, db.InsertActivityRoleParams{
		GuildID:        guildId,
		GrantType:      activityType,
		RoleID:         roleId,
		RequiredPoints: requiredPoints,
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, u.ErrActivityRoleExists
		}

		return nil, err
	}

	return nil, nil
}

func (uc *GuildUsecase) DeleteActivityRole(ctx context.Context, guildId string, roleId string) error {
	err := uc.q.DeleteActivityRole(ctx, db.DeleteActivityRoleParams{
		GuildID: guildId,
		RoleID:  roleId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (uc *GuildUsecase) GenerateGuildActivityLeaderboardCard(ctx context.Context, guildId string, acitivtyType, timePeriod string, page int) (gomponents.Node, error) {
	guild, err := uc.d.Guild(ctx, guildId)
	if err != nil {
		return nil, err
	}

	serverInfo := layouts.ServerInfo{
		Icon: guild.IconURL("100"),
		Name: guild.Name,
	}

	limitBy := int32(15)
	var card gomponents.Node
	switch timePeriod {
	case "weekly":
		leaderboard, err := uc.q.GetWeeklyActivityLeaderboard(ctx, db.GetWeeklyActivityLeaderboardParams{
			GuildID:   guildId,
			GrantType: acitivtyType,
			OffsetBy:  int32(page-1) * limitBy,
		})

		if err != nil {
			return nil, err
		}

		userIds := make([]string, 0)
		for _, value := range leaderboard {
			userIds = append(userIds, value.MemberID)
		}
		err = uc.d.RequestGuildMembersList(ctx, guildId, userIds, 0, "", true)
		if err != nil {
			return nil, err
		}

		fields := make([]layouts.LeaderboardDataField, 0)
		for _, value := range leaderboard {
			member, err := uc.d.GuildMember(ctx, guildId, value.MemberID)

			if member == nil || err != nil {
				fields = append(fields, layouts.LeaderboardDataField{
					Rank:     int(value.Rank),
					Username: value.MemberID,
					Value:    int(value.EarnedPoints),
				})
				continue
			}

			fields = append(fields, layouts.LeaderboardDataField{
				Rank:     int(value.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(value.EarnedPoints),
			})
		}

		card = layouts.ServerLeaderboard(layouts.ServerLeaderboardProps{
			ServerInfo: serverInfo,
			LeaderboardInfo: layouts.LeaderboardInfo{
				Name: "Activity Points - Weekly",
				Data: fields,
			},
		})
	case "monthly":
		leaderboard, err := uc.q.GetMonthlyActivityLeaderboard(ctx, db.GetMonthlyActivityLeaderboardParams{
			GuildID:   guildId,
			GrantType: acitivtyType,
			OffsetBy:  int32(page-1) * limitBy,
		})

		if err != nil {
			return nil, err
		}

		userIds := make([]string, 0)
		for _, value := range leaderboard {
			userIds = append(userIds, value.MemberID)
		}
		err = uc.d.RequestGuildMembersList(ctx, guildId, userIds, 0, "", true)
		if err != nil {
			return nil, err
		}

		fields := make([]layouts.LeaderboardDataField, 0)
		for _, value := range leaderboard {
			member, err := uc.d.GuildMember(ctx, guildId, value.MemberID)

			if member == nil || err != nil {
				fields = append(fields, layouts.LeaderboardDataField{
					Rank:     int(value.Rank),
					Username: value.MemberID,
					Value:    int(value.EarnedPoints),
				})
				continue
			}

			fields = append(fields, layouts.LeaderboardDataField{
				Rank:     int(value.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(value.EarnedPoints),
			})
		}

		card = layouts.ServerLeaderboard(layouts.ServerLeaderboardProps{
			ServerInfo: serverInfo,
			LeaderboardInfo: layouts.LeaderboardInfo{
				Name: "Activity Points - Monthly",
				Data: fields,
			},
		})
	default:
		leaderboard, err := uc.q.GetAllTimeActivityLeaderboard(ctx, db.GetAllTimeActivityLeaderboardParams{
			ActivityType: acitivtyType,
			GuildID:      guildId,
			LimitBy:      limitBy,
			OffsetBy:     int32(page-1) * limitBy,
		})

		if err != nil {
			return nil, err
		}

		userIds := make([]string, 0)
		for _, value := range leaderboard {
			userIds = append(userIds, value.MemberID)
		}
		err = uc.d.RequestGuildMembersList(ctx, guildId, userIds, 0, "", true)
		if err != nil {
			return nil, err
		}

		fields := make([]layouts.LeaderboardDataField, 0)
		for _, value := range leaderboard {
			member, err := uc.d.GuildMember(ctx, guildId, value.MemberID)

			if member == nil || err != nil {
				fields = append(fields, layouts.LeaderboardDataField{
					Rank:     int(value.Rank),
					Username: value.MemberID,
					Value:    int(value.Points),
				})
				continue
			}

			fields = append(fields, layouts.LeaderboardDataField{
				Rank:     int(value.Rank),
				Username: fmt.Sprintf("@%s", member.User.Username),
				Value:    int(value.Points),
			})
		}

		card = layouts.ServerLeaderboard(layouts.ServerLeaderboardProps{
			ServerInfo: serverInfo,
			LeaderboardInfo: layouts.LeaderboardInfo{
				Name: "Activity Points - All Time",
				Data: fields,
			},
		})
	}

	return card, nil
}

func (uc *GuildUsecase) CreateVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string, settings u.VoiceRoomLobbySettings) error {
	_, err := uc.q.CreateVoiceRoomLobby(ctx, db.CreateVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: originChannelId,

		UserLimit:      sqlx.Int32(settings.UserLimit),
		CanRename:      sqlx.Bool(settings.CanRename),
		CanLock:        sqlx.Bool(settings.CanLock),
		CanAdjustLimit: sqlx.Bool(settings.CanAdjustLimit),
	})

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return u.ErrVoiceRoomLobbyExists
		}

		return err
	}

	return err
}

func (uc *GuildUsecase) UpdateVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string, settings u.VoiceRoomLobbySettings) error {
	err := uc.q.UpdateVoiceRoomLobby(ctx, db.UpdateVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: originChannelId,

		UserLimit:      sqlx.Int32(settings.UserLimit),
		CanRename:      sqlx.Bool(settings.CanRename),
		CanLock:        sqlx.Bool(settings.CanLock),
		CanAdjustLimit: sqlx.Bool(settings.CanAdjustLimit),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u.ErrVoiceRoomLobbyNotFound
		}

		return err
	}

	return nil
}

func (uc *GuildUsecase) DeleteVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string) error {
	_, err := uc.q.GetVoiceRoomLobby(ctx, db.GetVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: originChannelId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u.ErrVoiceRoomLobbyNotFound
		}

		return err
	}

	err = uc.q.DeleteVoiceRoomLobby(ctx, db.DeleteVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: originChannelId,
	})

	if err != nil {
		return err
	}

	return nil
}

func (uc *GuildUsecase) RegisterVoiceRoom(ctx context.Context, guildId string, originChannelId string, channelId string, creatorUserId string) (*u.VoiceRoom, error) {
	room, err := uc.q.RegisterVoiceRoom(ctx, db.RegisterVoiceRoomParams{
		GuildID:         guildId,
		OriginChannelID: originChannelId,
		ChannelID:       channelId,
		CreatedByUserID: creatorUserId,
		CurrentOwnerID:  creatorUserId,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, u.ErrVoiceRoomExists
		}

		return nil, err
	}

	settings, err := uc.q.GetVoiceRoomLobby(ctx, db.GetVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: room.OriginChannelID,
	})
	if err != nil {
		return nil, err
	}

	return &u.VoiceRoom{
		OriginChannelId: room.OriginChannelID,
		CreatorId:       room.CreatedByUserID,
		CurrentOwnerId:  room.CurrentOwnerID,
		IsLocked:        room.IsLocked.Valid && room.IsLocked.Bool,

		Settings: u.VoiceRoomLobbySettings{
			UserLimit:      &settings.UserLimit,
			CanRename:      &settings.CanRename,
			CanLock:        &settings.CanLock,
			CanAdjustLimit: &settings.CanAdjustLimit,
		},
	}, nil
}

func (uc *GuildUsecase) GetVoiceRoom(ctx context.Context, guildId string, channelId string) (*u.VoiceRoom, error) {
	room, err := uc.q.GetVoiceRoom(ctx, db.GetVoiceRoomParams{
		GuildID:   guildId,
		ChannelID: channelId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, u.ErrVoiceRoomNotFound
		}

		return nil, err
	}

	settings, err := uc.q.GetVoiceRoomLobby(ctx, db.GetVoiceRoomLobbyParams{
		GuildID:        guildId,
		VoiceChannelID: room.OriginChannelID,
	})
	if err != nil {
		return nil, err
	}

	return &u.VoiceRoom{
		OriginChannelId: room.OriginChannelID,
		CreatorId:       room.CreatedByUserID,
		CurrentOwnerId:  room.CurrentOwnerID,
		IsLocked:        room.IsLocked.Valid && room.IsLocked.Bool,

		Settings: u.VoiceRoomLobbySettings{
			UserLimit:      &settings.UserLimit,
			CanRename:      &settings.CanRename,
			CanLock:        &settings.CanLock,
			CanAdjustLimit: &settings.CanAdjustLimit,
		},
	}, nil
}

func (uc *GuildUsecase) UpdateVoiceRoom(ctx context.Context, guildId string, channelId string, opts u.VoiceRoomModify) error {
	_, err := uc.q.UpdateVoiceRoom(ctx, db.UpdateVoiceRoomParams{
		GuildID:   guildId,
		ChannelID: channelId,

		CurrentOwnerID: sqlx.String(opts.CurrentOwnerId),
		IsLocked:       sqlx.Bool(opts.IsLocked),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u.ErrVoiceRoomNotFound
		}

		return err
	}

	return nil
}

func (uc *GuildUsecase) DeleteVoiceRoom(ctx context.Context, guildId string, channelId string) error {
	err := uc.q.DeleteVoiceRoom(ctx, db.DeleteVoiceRoomParams{
		GuildID:   guildId,
		ChannelID: channelId,
	})

	if err != nil {
		return err
	}

	return nil
}
