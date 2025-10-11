package discord_state

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"github.com/typical-developers/discord-bot-backend/pkg/redisx"
)

func (s *StateManager) GuildRoles(ctx context.Context, guildId string) ([]*discordgo.Role, error) {
	key := fmt.Sprintf("guild:%s:roles", guildId)

	result, err, _ := s.sf.Do(fmt.Sprintf("%s_roles", guildId), func() (any, error) {
		var roles []*discordgo.Role
		err := redisx.JSONUnwrap(ctx, s.redis, key, "$", &roles)

		if err != nil {
			if err != redis.Nil {
				return nil, err
			}

			roles, err := s.Session.GuildRoles(guildId, discordgo.WithContext(ctx), discordgo.WithRetryOnRatelimit(true))
			if err != nil {
				return nil, err
			}

			pipeline := s.redis.Pipeline()
			for _, role := range roles {
				key := fmt.Sprintf("guild:%s:role:%s", guildId, role.ID)
				pipeline.JSONSet(ctx, key, "$", role)
				pipeline.Expire(ctx, key, time.Hour)
			}

			_, _ = pipeline.Exec(ctx)

			return roles, nil
		}

		return roles, nil
	})

	if err != nil {
		return nil, err
	}

	return result.([]*discordgo.Role), nil
}

func (s *StateManager) GuildRole(ctx context.Context, guildId, roleId string) (*discordgo.Role, error) {
	key := fmt.Sprintf("guild:%s:role:%s", guildId, roleId)

	result, err, _ := s.sf.Do(fmt.Sprintf("%s_role:%s", guildId, roleId), func() (any, error) {
		var role *discordgo.Role
		err := redisx.JSONUnwrap(ctx, s.redis, key, "$", &role)

		if err != nil {
			if err != redis.Nil {
				return nil, err
			}

			roles, err := s.GuildRoles(ctx, guildId)
			if err != nil {
				return nil, err
			}

			for _, r := range roles {
				if r.ID == roleId {
					return r, nil
				}
			}

			return nil, ErrRoleNotFound
		}

		return role, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*discordgo.Role), nil
}
