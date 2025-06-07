package cache

import (
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

var (
	Client *redis.Client
)

func init() {
	var opts redis.Options

	if host := os.Getenv("REDIS_HOST"); host != "" {
		opts.Addr = host
	}

	if password := os.Getenv("REDIS_PASSWORD"); password != "" {
		opts.Password = password
	}

	if db := os.Getenv("REDIS_DB"); db != "" {
		dbNum, err := strconv.Atoi(db)
		if err != nil {
			logger.Log.WithSource.Error("Failed to parse REDIS_DB as int, using default (0) instead.", "error", err)
		} else {
			opts.DB = dbNum
		}
	}

	Client = redis.NewClient(&opts)
}
