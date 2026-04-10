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
		return kw
	case "vector":
		return vec
	default:
		return 0.65*kw + 0.35*vec
	}
}

func tokens(s string) map[string]int {
	out := map[string]int{}
	for _, w := range strings.Fields(strings.ToLower(s)) {
		if len(w) <= 2 {
			continue
		}
		out[w]++
	}
	return out
}

func keywordScore(query, content string) float64 {
	q := tokens(query)
	if len(q) == 0 {
		return 0
	}
	c := tokens(content)
	hits := 0
	for w := range q {
		if n, ok := c[w]; ok && n > 0 {
			hits++
		}
	}
	return float64(hits) / float64(len(q))
}

func vector(query, content string) []float64 {
	vocab := map[string]int{}
	q := tokens(query)
	c := tokens(content)
	for w := range q {
		if _, ok := vocab[w]; !ok {
			vocab[w] = len(vocab)
		}
	}
	for w := range c {
		if _, ok := vocab[w]; !ok {
			vocab[w] = len(vocab)
		}
	}
	v := make([]float64, len(vocab))
	for w, n := range q {
		if i, ok := vocab[w]; ok {
			v[i] = float64(n)
		}
	}
	return v
}

func cosineSimilarity(query, content string) float64 {
	q := tokens(query)
	c := tokens(content)
	if len(q) == 0 || len(c) == 0 {
		return 0
	}
	dot, nq, nc := 0.0, 0.0, 0.0
	for w, n := range q {
		ncnt, ok := c[w]
		_ = ncnt
		if ok {
			dot += float64(n * c[w])
		}
		nq += float64(n * n)
	}
	for _, n := range c {
		nc += float64(n * n)
	}
	if nq == 0 || nc == 0 {
		return 0
	}
	return dot / (math.Sqrt(nq) * math.Sqrt(nc))
}
