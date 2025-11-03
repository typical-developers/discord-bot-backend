package discord_state

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/typical-developers/discord-bot-backend/pkg/redisx"
)

func (s *StateManager) User(ctx context.Context, userId string) (*discordgo.User, error) {
	key := fmt.Sprintf("user:%s", userId)

	result, err, _ := s.sf.Do(fmt.Sprintf("%s_user:%s", userId, userId), func() (any, error) {
		var user *discordgo.User
		err := redisx.JSONUnwrap(ctx, s.redis, key, "$", &user)

		if err != nil {
			if err != redis.Nil {
				return nil, err
			}

			user, err := s.Session.User(userId, discordgo.WithContext(ctx), discordgo.WithRetryOnRatelimit(true))
			if err != nil {
				return nil, err
			}

			pipeline := s.redis.Pipeline()
			pipeline.JSONSet(ctx, key, "$", user)
			pipeline.Expire(ctx, key, time.Hour*24)
			_, _ = pipeline.Exec(ctx)

			return user, nil
		}

		return user, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*discordgo.User), nil
}
