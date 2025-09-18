package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// The key used to authorize access to the API.
	AuthKey string `env:"AUTH_KEY,required"`

	// The token used to authorize the Discord bot.
	DiscordToken string `env:"DISCORD_TOKEN,required"`

	// A PostgreSQL instance used to store data for the bot.
	//
	// Options are query parameters used in the connection string.
	// They should be formatted as an array, for example: "sslmode=disable,timezone=utc".
	Database struct {
		Host     string   `env:"HOST,required"`
		Password string   `env:"PASSWORD"`
		Port     int      `env:"PORT,required"`
		Options  []string `env:"OPTIONS"`
	} `envPrefix:"DATABASE_"`

	// A Redis instance used to store data fetched from the PostgreSQL database instance.
	DatabaseCache struct {
		Host     string `env:"HOST,required"`
		Password string `env:"PASSWORD"`
		Port     int    `env:"PORT,required"`
		DB       int    `env:"DB,required"`
	} `envPrefix:"DATABASE_CACHE_"`

	// A Redis instance  used to cache Discord API responses.
	DiscordCache struct {
		Host     string `env:"HOST,required"`
		Password string `env:"PASSWORD"`
		Port     int    `env:"PORT,required"`
		DB       int    `env:"DB,required"`
	} `envPrefix:"DISCORD_CACHE_"`
}

var (
	C Config

	once sync.Once
)

func init() {
	once.Do(func() {
		if err := env.Parse(&C); err != nil {
			panic(err)
		}
	})
}
