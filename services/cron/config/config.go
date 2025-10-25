package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// The port to run the API on.
	Port int `env:"PORT" envDefault:"8080"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`

	// A PostgreSQL instance used to store data for the bot.
	//
	// Options are query parameters used in the connection string.
	// They should be formatted as an array, for example: "sslmode=disable,timezone=utc".
	Database struct {
		Username string `env:"USERNAME,required"`
		Password string `env:"PASSWORD"`
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT,required"`
		Options  string `env:"OPTIONS"`
	} `envPrefix:"DATABASE_"`
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
