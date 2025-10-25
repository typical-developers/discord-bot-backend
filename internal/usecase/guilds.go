package usecase

import (
	"context"

	"maragu.dev/gomponents"
)

type GuildsUsecase interface {
	RegisterGuild(ctx context.Context, guildId string) (*GuildSettings, error)
	GetGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)

	UpdateGuildActivitySettings(ctx context.Context, guildId string, opts UpdateAcitivtySettings) (*GuildSettings, error)

	CreateActivityRole(ctx context.Context, guildId string, activityType string, roleId string, requiredPoints int32) (*GuildActivityRole, error)
	DeleteActivityRole(ctx context.Context, guildId string, roleId string) error

	GenerateGuildActivityLeaderboardCard(ctx context.Context, guildId string, acitivtyType, timePeriod string, page int) (gomponents.Node, error)

	CreateVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string, settings VoiceRoomLobbySettings) error
	UpdateVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string, settings VoiceRoomLobbySettings) error
	DeleteVoiceRoomLobby(ctx context.Context, guildId string, originChannelId string) error
	RegisterVoiceRoom(ctx context.Context, guildId string, originChannelId string, channelId string, creatorUserId string) (*VoiceRoom, error)
	GetVoiceRoom(ctx context.Context, guildId string, channelId string) (*VoiceRoom, error)
	UpdateVoiceRoom(ctx context.Context, guildId string, channelId string, opts VoiceRoomModify) error
	DeleteVoiceRoom(ctx context.Context, guildId string, channelId string) error
}
