package discord_state

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

const ()

type StateManager struct {
	Session *discordgo.Session
	redis   *redis.Client

	sf singleflight.Group
}

type StateManagerOptions struct {
	// The DiscordGo instance to use for interacting.
	DiscordSession *discordgo.Session

	// The redis client instance to use for caching.
	RedisClient *redis.Client
}

func NewStateManager(opts *StateManagerOptions) *StateManager {
	state := &StateManager{
		Session: opts.DiscordSession,
		redis:   opts.RedisClient,
	}

	opts.DiscordSession.AddHandler(func(s *discordgo.Session, e any) {
		ctx := context.Background()

		switch e := e.(type) {
		case *discordgo.GuildUpdate:
			state.redis.JSONSet(ctx, fmt.Sprintf("guild:%s", e.ID), "$", e.Guild)
		case *discordgo.GuildDelete:
			state.redis.JSONDel(ctx, fmt.Sprintf("guild:%s", e.ID), "$")
		case *discordgo.GuildMemberUpdate:
			state.redis.JSONSet(ctx, fmt.Sprintf("guild:%s:member:%s", e.GuildID, e.User.ID), "$", e.Member)
		case *discordgo.GuildMembersChunk:
			for _, member := range e.Members {
				state.redis.JSONSet(ctx, fmt.Sprintf("guild:%s:member:%s", e.GuildID, member.User.ID), "$", member)
			}
		case *discordgo.GuildMemberRemove:
			state.redis.JSONDel(ctx, fmt.Sprintf("guild:%s:member:%s", e.GuildID, e.User.ID), "$")
		case *discordgo.GuildRoleCreate:
			state.redis.JSONSet(ctx, fmt.Sprintf("guild:%s:role:%s", e.GuildID, e.Role.ID), "$", e.Role)
		case *discordgo.GuildRoleUpdate:
			state.redis.JSONSet(ctx, fmt.Sprintf("guild:%s:role:%s", e.GuildID, e.Role.ID), "$", e.Role)
		case *discordgo.GuildRoleDelete:
			state.redis.JSONDel(ctx, fmt.Sprintf("guild:%s:role:%s", e.GuildID, e.RoleID), "$")
		}
	})

	return state
}
