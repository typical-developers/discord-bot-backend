package db

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

var (
	once sync.Once
	Pool *pgxpool.Pool
)

const (
	MinConnections = 0
	MaxConnections = 10
)

func InitalizePool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	connectionUrl := os.Getenv("POSTGRES_URL")

	if connectionUrl == "" {
		return nil, fmt.Errorf("POSTGRES_URL is not set.")
	}

	var dbError error
	once.Do(func() {
		var config *pgxpool.Config

		config, dbError = pgxpool.ParseConfig(connectionUrl)
		config.MinConns = MinConnections
		config.MaxConns = MaxConnections

		config.BeforeAcquire = func(ctx context.Context, pgConn *pgx.Conn) bool {
			logger.Log.Debug("Acquiring a connection to the database pool.")
			return true
		}

		config.AfterRelease = func(pgConn *pgx.Conn) bool {
			logger.Log.Debug("Released a connection back into the pool.")
			return true
		}

		Pool, dbError = pgxpool.NewWithConfig(ctx, config)

		if dbError != nil {
			return
		}

		dbError = Pool.Ping(ctx)
		if dbError != nil {
			return
		}
	})

	return Pool, dbError
}

func init() {
	_, err := InitalizePool()
	if err != nil {
		panic(err)
	}
}
