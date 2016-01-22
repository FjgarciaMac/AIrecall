// Package storage owns the SQLite schema for episodes and facts.
package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps the sqlite connection.
type Store struct {
	DB *sql.DB
}

// Open opens (and migrates) the store.
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS episodes (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			agent_id   TEXT NOT NULL,
			content    TEXT NOT NULL,
			kind       TEXT NOT NULL DEFAULT 'episode',
			created_at TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS facts (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			agent_id   TEXT NOT NULL,
			key        TEXT NOT NULL,
			value      TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			UNIQUE(agent_id, key)
		);
		CREATE INDEX IF NOT EXISTS idx_episodes_agent
			ON episodes(agent_id, created_at);
	`); err != nil {
		db.Close()
		return nil, err
