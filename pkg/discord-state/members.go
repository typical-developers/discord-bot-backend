package discord_state

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/typical-developers/discord-bot-backend/pkg/redisx"
)

func (s *StateManager) GuildMember(ctx context.Context, guildId, userId string) (*discordgo.Member, error) {
	key := fmt.Sprintf("guild:%s:member:%s", guildId, userId)

	result, err, _ := s.sf.Do(fmt.Sprintf("%s_member:%s", guildId, userId), func() (any, error) {
		var member *discordgo.Member
		err := redisx.JSONUnwrap(ctx, s.redis, key, "$", &member)

		if err != nil {
			if err != redis.Nil {
				return nil, err
			}

			member, err := s.Session.GuildMember(guildId, userId, discordgo.WithContext(ctx), discordgo.WithRetryOnRatelimit(true))
			if err != nil {
				return nil, err
			}

			pipeline := s.redis.Pipeline()
			pipeline.JSONSet(ctx, key, "$", member)
			_, _ = pipeline.Exec(ctx)

			return member, nil
		}

		return member, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*discordgo.Member), nil
}
