package tasks

import (
	"database/sql"

	"github.com/typical-developers/discord-bot-backend/internal/db"
)

type Tasks struct {
	db *sql.DB
	q  *db.Queries
}

func NewTasks(db *sql.DB, q *db.Queries) *Tasks {
	return &Tasks{db: db, q: q}
}
