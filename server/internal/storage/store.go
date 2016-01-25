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
	}
	return &Store{DB: db}, nil
}

// Close releases the connection.
func (s *Store) Close() error { return s.DB.Close() }

// AddEpisode stores one episodic memory.
func (s *Store) AddEpisode(agentID, content string) error {
	_, err := s.DB.Exec(
		"INSERT INTO episodes (agent_id, content, created_at) VALUES (?, ?, ?)",
		agentID, content, time.Now().UTC().Format(time.RFC3339))
	return err
}

// UpsertFact stores or updates a semantic fact.
func (s *Store) UpsertFact(agentID, key, value string) error {
	_, err := s.DB.Exec(`
		INSERT INTO facts (agent_id, key, value, updated_at) VALUES (?, ?, ?, ?)
