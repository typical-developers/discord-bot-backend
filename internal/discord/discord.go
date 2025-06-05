package discord

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/typical-developers/discord-bot-backend/pkg/redisutil"
)

type Cache struct {
}

type DiscordClient struct {
	*discordgo.Session
	Cache *Cache
}

var (
	Client *DiscordClient
)

func (c *Cache) Guild(guildId string) (*discordgo.Guild, error) {
	ctx := context.Background()
	key := fmt.Sprintf("guild_%s", guildId)

	var guildInfo *discordgo.Guild
	if cache := redisutil.GetCached[discordgo.Guild](ctx, key); cache != nil {
		guildInfo = cache
	} else {
		guild, err := Client.Guild(guildId)
		if err != nil {
			return nil, err
		}

		guildInfo = guild

		duration := 60 * time.Minute
		redisutil.SetCached(ctx, key, guild, &redisutil.CacheOpts{
			Expiry: &duration,
		})
	}

	return guildInfo, nil
}

func (c *Cache) GuildMember(guildId, memberId string) (*discordgo.Member, error) {
	ctx := context.Background()
	key := fmt.Sprintf("guild_%s:member_%s", guildId, memberId)

	var memberInfo *discordgo.Member
	if cache := redisutil.GetCached[discordgo.Member](ctx, key); cache != nil {
		memberInfo = cache
	} else {
		member, err := Client.GuildMember(guildId, memberId)
		if err != nil {
			return nil, err
		}

		memberInfo = member

		duration := 60 * time.Minute
		redisutil.SetCached(context.Background(), key, member, &redisutil.CacheOpts{
			Expiry: &duration,
		})
	}

	return memberInfo, nil
}

func init() {
	c, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	c.StateEnabled = true

	Client = &DiscordClient{
		Session: c,
		Cache:   &Cache{},
	}
}
