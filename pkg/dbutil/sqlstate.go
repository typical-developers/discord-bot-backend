package dbutil

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Errors to help identify SQLState codes more cleanly.
type SQLState string

const (
	SQLStateUniqueViolation SQLState = "23505"
)

// Utility to unwrap a pgconn.PgError cleanly.
func UnwrapSQLState(err error) (SQLState, bool) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return SQLState(pgErr.Code), true
	}

	return "", false
}
