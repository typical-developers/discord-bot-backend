package tasks

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func (t *Tasks) FlushWeeklyActivityLeaderboard(ctx context.Context) error {
	details, err := t.q.GetWeeklyActivityLeaderboardResetDetails(ctx)
	if err != nil {
		return err
	}

	if details.LastReset >= details.ExpectedReset {
		log.WithFields(log.Fields{
			"row_total":      details.RowTotal,
			"last_reset":     details.LastReset,
			"expected_reset": details.ExpectedReset,
		}).Info("The weekly activity leaderboard has already been reset.")
		return nil
	}

	if details.RowTotal <= 0 {
		log.WithFields(log.Fields{
			"row_total":      details.RowTotal,
			"last_reset":     details.LastReset,
			"expected_reset": details.ExpectedReset,
		}).Info("The current weekly activity leaderboard is empty, no reset necessary.")
		return nil
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := t.q.WithTx(tx)
	if err := q.ArchiveWeeklyActivityLeaderboard(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := q.FlushOudatedWeeklyActivityLeaderboard(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"row_total":      details.RowTotal,
		"last_reset":     details.LastReset,
		"expected_reset": details.ExpectedReset,
	}).Info("The weekly activity leaderboard has been reset.")

	return nil
}

func (t *Tasks) FlushMonthlyActivityLeaderboard(ctx context.Context) error {
	details, err := t.q.GetMonthlyActivityLeaderboardResetDetails(ctx)
	if err != nil {
		return err
	}

	if details.LastReset >= details.ExpectedReset {
		log.WithFields(log.Fields{
			"row_total":      details.RowTotal,
			"last_reset":     details.LastReset,
			"expected_reset": details.ExpectedReset,
		}).Info("The current monthly activity leaderboard has already been reset.")
		return nil
	}

	if details.RowTotal <= 0 {
		log.WithFields(log.Fields{
			"row_total":      details.RowTotal,
			"last_reset":     details.LastReset,
			"expected_reset": details.ExpectedReset,
		}).Info("The current monthly activity leaderboard is empty, no reset necessary.")
		return nil
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := t.q.WithTx(tx)
	if err := q.ArchiveMonthlyActivityLeaderboard(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := q.FlushOudatedMonthlyActivityLeaderboard(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"row_total":      details.RowTotal,
		"last_reset":     details.LastReset,
		"expected_reset": details.ExpectedReset,
	}).Info("The monthly activity leaderboard has been reset")

	return nil
}
