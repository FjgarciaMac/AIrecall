// Package summarizer distills durable facts from older episodes.
package summarizer

import (
	"regexp"
	"strings"

	"github.com/FjgarciaMac/AIrecall/server/internal/storage"
)

// Extractor finds "key: value" style facts in episode text.
var factPattern = regexp.MustCompile(
	`(?i)\b(?:prefers|preferred|likes|wants|uses|contact|timezone|language|name|email|phone)\b[^,.;\n]{0,60}`)

// Summarizer runs the compaction pass.
type Summarizer struct {
	store *storage.Store
}

// New creates a summarizer bound to the store.
func New(store *storage.Store) *Summarizer { return &Summarizer{store: store} }

// Result reports what one Run pass promoted.
type Result struct {
	Scanned       int `json:"scanned"`
	FactsPromoted int `json:"facts_promoted"`
}

// Run walks the newest episodes and promotes detectable facts into the
// facts table. Compaction (deleting superseded episodes) is left to the
// server operator; this pass only distills.
func (s *Summarizer) Run(agentID string, limit int) (Result, error) {
	episodes, err := s.store.RecentEpisodes(agentID, limit)
	if err != nil {
		return Result{}, err
	}
	result := Result{Scanned: len(episodes)}
	for _, e := range episodes {
		kv := splitFact(e.Content)
		if kv == nil {
			continue
		}
		if err := s.store.UpsertFact(agentID, kv[0], kv[1]); err != nil {
			return Result{}, err
		}
		result.FactsPromoted++
	}
	return result, nil
}

// splitFact extracts a key/value pair from one episode line, or nil when
// the line carries no recognizable fact trigger. Accepted separators, in
// priority order: ":", "=", " is ", and a bare space.
func splitFact(text string) []string {
	m := factPattern.FindString(text)
	if m == "" {
		return nil
	}
