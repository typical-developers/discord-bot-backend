// All of these are used in a separate package that handles tasks.
//
// Tasks aren't on the backend itself in the event that the API itself goes down.
// Allows tasks to run still to keep cleanups and data updated.
//
// This should also make it easier to scale up if necessary.

package dbutil

import (
	"context"

	"github.com/typical-developers/discord-bot-backend/internal/db"
)

func ResetWeeklyActivityLeaderboards(ctx context.Context) error {
	connection, err := Client(ctx)
	if err != nil {
		return err
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		return err
	}
	queries := db.New(connection).WithTx(tx)
	defer connection.Release()

	err = queries.ResetWeeklyActivityLeaderboard(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = queries.TruncateWeeklyActivityLeaderboard(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ResetMonthlyActivityLeaderboards(ctx context.Context) error {
	connection, err := Client(ctx)
	if err != nil {
		return err
	}

	tx, err := connection.Begin(ctx)
	if err != nil {
		return err
	}
	queries := db.New(connection).WithTx(tx)
	defer connection.Release()

	err = queries.ResetMonthlyActivityLeaderboard(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = queries.TruncateMonthlyActivityLeaderboard(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
