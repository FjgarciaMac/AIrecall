// Package retrieval implements hybrid search: BM25-style keyword scoring
// combined with cosine similarity over word-overlap vectors.
package retrieval

import (
	"math"
	"sort"
	"strings"
)

// Memory is the retrieval unit (mirrors storage.Episode).
type Memory struct {
	Content string
}

// Scorer scores one query against one memory.
type Scorer struct{}

// NewScorer returns a default scorer.
func NewScorer() *Scorer { return &Scorer{} }

// Hybrid ranks memories by combined keyword + vector score.
// mode: "keyword" | "vector" | "hybrid"
func (s *Scorer) Hybrid(query string, memories []Memory, topK int, mode string) []string {
	type scored struct {
		m  Memory
		sc float64
	}
	results := make([]scored, 0, len(memories))
	for _, m := range memories {
		sc := s.Score(query, m.Content, mode)
		if sc > 0 {
			results = append(results, scored{m, sc})
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].sc > results[j].sc })
	if len(results) > topK {
		results = results[:topK]
	}
	out := make([]string, 0, len(results))
	for _, r := range results {
		out = append(out, r.m.Content)
	}
	return out
}

// Score combines keyword and vector similarity.
func (s *Scorer) Score(query, content, mode string) float64 {
	kw := keywordScore(query, content)
	vec := cosineSimilarity(query, content)
	switch mode {
	case "keyword":
