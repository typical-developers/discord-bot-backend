package usecase

import "context"

type GuildsUsecase interface {
	CreateGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)
	GetGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)
}
