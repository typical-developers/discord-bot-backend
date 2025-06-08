package tasks

import (
	"context"

	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func FlushWeeklyActivityLeaderboard() {
	ctx := context.Background()
	err := dbutil.ResetWeeklyActivityLeaderboards(ctx)

	if err != nil {
		logger.Log.Error("Failed to reset weekly activity leaderboards.", "error", err)
		return
	}

	logger.Log.Info("Weekly activity leaderboards successfully reset.")
}
