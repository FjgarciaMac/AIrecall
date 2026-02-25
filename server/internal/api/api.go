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
	if !ok {
		return
	}
	agent := str(body, "agent_id", "default")
	query := str(body, "query", "")
	topK := intF(body, "top_k", float64(s.cfg.TopK))
	if topK <= 0 {
		topK = s.cfg.TopK
	}
	episodes, err := s.store.RecentEpisodes(agent, s.cfg.RecallLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	memories := make([]retrieval.Memory, 0, len(episodes))
	for _, e := range episodes {
		memories = append(memories, retrieval.Memory{Content: e.Content})
	}
	hits := s.scorer.Hybrid(query, memories, topK, s.cfg.Mode)
	writeJSON(w, map[string]any{"memories": hits})
}

func (s *Server) handleGetFact(w http.ResponseWriter, r *http.Request) {
	agent := r.PathValue("agent")
	key := r.PathValue("key")
	value, ok, err := s.store.Fact(agent, key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, map[string]any{"key": key, "value": value})
}

func (s *Server) handleSummarize(w http.ResponseWriter, r *http.Request) {
	body, ok := s.decode(w, r)
	if !ok {
		return
	}
	agent := str(body, "agent_id", "default")
	result, err := s.summarizer.Run(agent, s.cfg.RecallLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"facts_promoted": result.FactsPromoted,
	})
}

func str(m map[string]any, key, fallback string) string {
	if v, ok := m[key].(string); ok && v != "" {
		return v
	}
	return fallback
}

func intF(m map[string]any, key string, fallback float64) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return int(fallback)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
// draft note 1412
