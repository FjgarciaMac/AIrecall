package retrieval

import "testing"

func TestKeywordScoreExact(t *testing.T) {
	s := NewScorer()
	sc := s.Score("refund policy", "User asked about refund policy", "keyword")
	if sc != 1.0 {
		t.Fatalf("expected 1.0, got %f", sc)
	}
