package usecase

import "context"

type GuildsUsecase interface {
	CreateGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)
	GetGuildSettings(ctx context.Context, guildId string) (*GuildSettings, error)
	UpdateGuildActivitySettings(ctx context.Context, guildId string, opts UpdateAcitivtySettings) (*GuildSettings, error)
}
