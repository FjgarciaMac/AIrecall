// Package api wires the HTTP handlers for the memory server.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/FjgarciaMac/AIrecall/server/internal/retrieval"
	"github.com/FjgarciaMac/AIrecall/server/internal/storage"
	"github.com/FjgarciaMac/AIrecall/server/internal/summarizer"
)

// Config holds server tuning.
type Config struct {
	TopK        int    `json:"top_k"`
	Mode        string `json:"mode"` // keyword | vector | hybrid
	RecallLimit int    `json:"recall_limit"`
}

// DefaultConfig returns the built-in defaults.
func DefaultConfig() Config {
	return Config{TopK: 5, Mode: "hybrid", RecallLimit: 200}
}

// Server holds the handlers.
type Server struct {
	store      *storage.Store
	scorer     *retrieval.Scorer
	summarizer *summarizer.Summarizer
	cfg        Config
}

// New builds the handler mux.
func New(store *storage.Store, cfg Config) http.Handler {
	s := &Server{
		store:      store,
		scorer:     retrieval.NewScorer(),
		summarizer: summarizer.New(store),
		cfg:        cfg,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/remember", s.handleRemember)
	mux.HandleFunc("POST /v1/facts", s.handleUpsertFact)
	mux.HandleFunc("POST /v1/recall", s.handleRecall)
	mux.HandleFunc("GET /v1/facts/{agent}/{key}", s.handleGetFact)
