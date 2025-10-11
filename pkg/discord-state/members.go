package discord_state

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

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
				var dgError *discordgo.RESTError
				if errors.As(err, &dgError) && dgError.Message.Code == discordgo.ErrCodeUnknownMember {
					// cache the value as nil, if they ever rejoin the guild it will automatically fetch their information.
					pipeline := s.redis.Pipeline()
					pipeline.JSONSet(ctx, key, "$", nil)
					pipeline.Expire(ctx, key, time.Hour*24)
					_, _ = pipeline.Exec(ctx)
				}

				return nil, err
			}

			pipeline := s.redis.Pipeline()
			pipeline.JSONSet(ctx, key, "$", member)
			pipeline.Expire(ctx, key, time.Hour*24)
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

func (s *StateManager) RequestGuildMembersList(ctx context.Context, guildId string, userIds []string, limit int, nonce string, presences bool) error {
	// This is to build the member list cache.
	//
	// If all the members are cached, it wont request the resources for them again.
	// Any resource that isn't cached will be requested.
	//
	// I'm not 100% certain if this is the most efficient way to handle this.
	// This relies on fetching a chunk then caching all of them separately to be used after the fact.

	_, err, _ := s.sf.Do(fmt.Sprintf("%s_members_list_%s", guildId, strings.Join(userIds, ",")), func() (any, error) {
		pipeline := s.redis.Pipeline()

		pipelineCmds := map[string]*redis.IntCmd{}
		for _, userId := range userIds {
			pipelineCmds[userId] = pipeline.Exists(ctx, fmt.Sprintf("guild:%s:member:%s", guildId, userId))
		}
		if _, err := pipeline.Exec(ctx); err != nil {
			return nil, err
		}

		for userId, cmd := range pipelineCmds {
			if cmd.Val() == 0 {
				continue
			}

			i := slices.Index(userIds, userId)
			userIds = slices.Delete(userIds, i, i+1)
		}

		if len(userIds) > 0 {
			err := s.Session.RequestGuildMembersList(guildId, userIds, limit, nonce, presences)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}
