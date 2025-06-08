package tasks

import (
	"context"

	"github.com/typical-developers/discord-bot-backend/pkg/dbutil"
	"github.com/typical-developers/discord-bot-backend/pkg/logger"
)

func FlushMonthlyActivityLeaderboard() {
	ctx := context.Background()
	err := dbutil.ResetMonthlyActivityLeaderboards(ctx)

	if err != nil {
		logger.Log.Error("Failed to reset monthly activity leaderboards.", "error", err)
		return
	}

	logger.Log.Info("Monthly activity leaderboards successfully reset.")
}
