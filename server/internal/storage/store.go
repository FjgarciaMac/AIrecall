// Package storage owns the SQLite schema for episodes and facts.
package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store wraps the sqlite connection.
