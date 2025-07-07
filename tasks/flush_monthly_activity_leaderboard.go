package tasks

import (
	"context"

	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func FlushMonthlyActivityLeaderboard() {
	ctx := context.Background()

	connection, err := dbutil.Client(ctx)
	if err != nil {
		logger.Log.Error("Failed to get database connection.", "error", err)
		return
	}
	queries := db.New(connection)
	defer connection.Release()

	epoches, err := queries.GetMonthlyActivityLeaderboardLastReset(ctx)
	if err != nil {
		logger.Log.Error("Failed to get the monthly activity leaderboard last reset time.", "error", err)
		return
	}

	if epoches.MonthStart == epoches.ExpectedMonthStart {
		logger.Log.Info("Monthly activity leaderboards have already been reset.")
		return
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.Error("Failed to start transaction.", "error", err)
		return
	}
	txQueries := db.New(connection).WithTx(tx)

	if err := txQueries.ResetMonthlyActivityLeaderboard(ctx); err != nil {
		logger.Log.Error("Failed to reset monthly activity leaderboards.", "error", err)
		return
	}

	if err := txQueries.TruncateMonthlyActivityLeaderboard(ctx); err != nil {
		logger.Log.Error("Failed to reset monthly activity leaderboards.", "error", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.Error("Failed to commit transaction.", "error", err)
		return
	}

	logger.Log.Info("Monthly activity leaderboards successfully reset.")
}
