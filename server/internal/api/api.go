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
	mux.HandleFunc("POST /v1/summarize", s.handleSummarize)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	return mux
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request) (map[string]any, bool) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return nil, false
	}
	return body, true
}

func (s *Server) handleRemember(w http.ResponseWriter, r *http.Request) {
	body, ok := s.decode(w, r)
	if !ok {
		return
	}
	agent := str(body, "agent_id", "default")
	content := str(body, "content", "")
	if content == "" {
		http.Error(w, "content required", http.StatusBadRequest)
		return
	}
	if err := s.store.AddEpisode(agent, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleUpsertFact(w http.ResponseWriter, r *http.Request) {
	body, ok := s.decode(w, r)
	if !ok {
		return
	}
	agent := str(body, "agent_id", "default")
	key := str(body, "key", "")
	value := str(body, "value", "")
	if key == "" {
		http.Error(w, "key required", http.StatusBadRequest)
		return
	}
	if err := s.store.UpsertFact(agent, key, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleRecall(w http.ResponseWriter, r *http.Request) {
	body, ok := s.decode(w, r)
