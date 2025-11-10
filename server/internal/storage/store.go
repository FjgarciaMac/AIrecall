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
		ON CONFLICT(agent_id, key) DO UPDATE SET value=excluded.value,
			updated_at=excluded.updated_at`,
		agentID, key, value, time.Now().UTC().Format(time.RFC3339))
	return err
}

// Fact returns a semantic fact by key.
func (s *Store) Fact(agentID, key string) (string, bool, error) {
	var value string
	err := s.DB.QueryRow(
		"SELECT value FROM facts WHERE agent_id = ? AND key = ?",
		agentID, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return value, true, nil
}

// Episode is one episodic memory row as returned by RecentEpisodes.
type Episode struct {
	ID      int64
	Content string
}

// RecentEpisodes returns the newest episodes for one agent, newest first.
func (s *Store) RecentEpisodes(agentID string, limit int) ([]Episode, error) {
	rows, err := s.DB.Query(
		`SELECT id, content FROM episodes WHERE agent_id = ?
		 ORDER BY created_at DESC, id DESC LIMIT ?`, agentID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Episode
	for rows.Next() {
		var e Episode
		if err := rows.Scan(&e.ID, &e.Content); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// CountEpisodes counts the episodes stored for one agent.
func (s *Store) CountEpisodes(agentID string) (int, error) {
	var n int
	err := s.DB.QueryRow(
		"SELECT COUNT(*) FROM episodes WHERE agent_id = ?", agentID).Scan(&n)
	return n, err
}
