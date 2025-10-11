package usecase

import (
	"context"

	"maragu.dev/gomponents"
)

type GuildsUsecase interface {
	CreateGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)
	GetGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)

	UpdateGuildActivitySettings(ctx context.Context, guildId string, opts UpdateAcitivtySettings) (*GuildSettings, error)

	CreateActivityRole(ctx context.Context, guildId string, activityType string, roleId string, requiredPoints int32) (*GuildActivityRole, error)
	DeleteActivityRole(ctx context.Context, guildId string, roleId string) error

	GenerateGuildActivityLeaderboard(ctx context.Context, guildId string, acitivtyType, timePeriod string, page int) (gomponents.Node, error)
}
