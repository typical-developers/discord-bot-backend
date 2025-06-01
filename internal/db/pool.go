package db

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

var (
	once sync.Once
	pool *pgxpool.Pool
)

const (
	MinConnections = 0
	MaxConnections = 10
)

func InitalizePool() (*pgxpool.Pool, error) {
	ctx := context.Background()
	connectionUrl := os.Getenv("POSTGRES_URL")

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

		pool, dbError = pgxpool.NewWithConfig(ctx, config)

		if dbError != nil {
			return
		}

		dbError = pool.Ping(ctx)
		if dbError != nil {
			return
		}
	})

	return pool, dbError
}

func Client(ctx context.Context) (*pgxpool.Conn, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
