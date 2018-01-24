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
