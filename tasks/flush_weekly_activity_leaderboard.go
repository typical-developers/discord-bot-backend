package tasks

import (
	"context"

	"github.com/typical-developers/discord-bot-backend/internal/db"
	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func FlushWeeklyActivityLeaderboard() {
	ctx := context.Background()

	connection, err := dbutil.Client(ctx)
	if err != nil {
		logger.Log.Error("Failed to get database connection.", "error", err)
		return
	}
	queries := db.New(connection)
	defer connection.Release()

	epoches, err := queries.GetWeeklyActivityLeaderboardLastReset(ctx)
	if err != nil {
		logger.Log.Error("Failed to get the weekly activity leaderboard last reset time.", "error", err)
		return
	}

	if epoches.WeekStart == epoches.ExpectedWeekStart {
		logger.Log.Info("Weekly activity leaderboards have already been reset.")
		return
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		logger.Log.Error("Failed to start transaction.", "error", err)
		return
	}
	txQueries := db.New(connection).WithTx(tx)

	if err := txQueries.ResetWeeklyActivityLeaderboard(ctx); err != nil {
		logger.Log.Error("Failed to reset weekly activity leaderboards.", "error", err)
		return
	}

	if err := txQueries.TruncateWeeklyActivityLeaderboard(ctx); err != nil {
		logger.Log.Error("Failed to reset weekly activity leaderboards.", "error", err)
		return
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Log.Error("Failed to commit transaction.", "error", err)
		return
	}

	logger.Log.Info("Weekly activity leaderboards successfully reset.")
}
