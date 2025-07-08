package db

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

const (
	MinConnections = 0
	MaxConnections = 10

	// How often to ping the database to make sure it's still alive.
	PingInterval = 30 * time.Second

	// How many attempts should there be to retry the connection.
	// This is used for both connecting / reconnecting.
	RetryAttempts = 10
	// How often to wait between each retry attempt.
	RetryDelay = 5 * time.Second
)

var (
	mu sync.RWMutex

	Pool   *pgxpool.Pool
	Config *pgxpool.Config
)

func aliveCheck(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(PingInterval):
			err := Pool.Ping(ctx)
			if err != nil {
				logger.Log.Warn("Database has been disconnected. Attempting to reconnect.")

				err := connect(ctx)
				if err != nil {
					logger.Log.Error("Failed to reconnect to the database.", "error", err)
					os.Exit(1)
				}

				continue
			}

			logger.Log.Debug("Database connection is alive.")
		}
	}
}

func connect(ctx context.Context) error {
	var returnErr error

	for i := range RetryAttempts {
		// In the event that this error happens, it's likely something to do with the config.
		// So in this case, we just exit the program.
		//
		// Realistically shouldn't ever happen anyway unless something REALLY dumb is being done.
		pgpool, err := pgxpool.NewWithConfig(ctx, Config)
		if err != nil {
			logger.Log.Error("Failed to create a new connection pool.", "error", err)
			os.Exit(1)
		}

		err = pgpool.Ping(ctx)
		if err != nil {
			if i == RetryAttempts-1 {
				returnErr = err
				break
			}

			logger.Log.Error(
				"Failed to connect to the database.",
				"attempt", i+1,
				"remainingAttempts", RetryAttempts-i-1,
				"error", err,
			)
			time.Sleep(RetryDelay)
			continue
		}

		mu.Lock()
		Pool = pgpool
		mu.Unlock()

		logger.Log.Info("Successfully connected to the database.")
		return nil
	}

	return returnErr
}

func init() {
	ctx := context.Background()
	connectionUrl := os.Getenv("POSTGRES_URL")

	if connectionUrl == "" {
		logger.Log.Error("POSTGRES_URL has not been set.")
		os.Exit(1)
	}

	config, err := pgxpool.ParseConfig(connectionUrl)
	if err != nil {
		logger.Log.Error("Failed to parse connection URL.", "error", err)
		os.Exit(1)
	}

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

	Config = config

	err = connect(ctx)
	if err != nil {
		logger.Log.Error("Failed to connect to the database.", "error", err)
		os.Exit(1)
	}

	go aliveCheck(ctx)
}
