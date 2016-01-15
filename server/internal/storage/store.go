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
