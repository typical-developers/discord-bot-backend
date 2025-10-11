package discord_state

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/typical-developers/discord-bot-backend/pkg/redisx"
)

func (s *StateManager) Guild(ctx context.Context, guildId string) (*discordgo.Guild, error) {
	key := fmt.Sprintf("guild:%s", guildId)

	result, err, _ := s.sf.Do(guildId, func() (any, error) {
		var guild *discordgo.Guild
		err := redisx.JSONUnwrap(ctx, s.redis, key, "$", &guild)

		if err != nil {
			if err != redis.Nil {
				return nil, err
			}

			guild, err = s.Session.Guild(guildId, discordgo.WithContext(ctx), discordgo.WithRetryOnRatelimit(true))
			if err != nil {
				return nil, err
			}

			pipeline := s.redis.Pipeline()
			pipeline.JSONSet(ctx, key, "$", guild)
			_, _ = pipeline.Exec(ctx)

			return guild, nil
		}

		return guild, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*discordgo.Guild), nil
}
