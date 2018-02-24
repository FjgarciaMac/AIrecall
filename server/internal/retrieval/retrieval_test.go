package retrieval

import "testing"

func TestKeywordScoreExact(t *testing.T) {
	s := NewScorer()
	sc := s.Score("refund policy", "User asked about refund policy", "keyword")
	if sc != 1.0 {
		t.Fatalf("expected 1.0, got %f", sc)
	}
}

func TestKeywordScorePartial(t *testing.T) {
	s := NewScorer()
	sc := s.Score("refund policy for orders", "refund window is 30 days", "keyword")
	if sc <= 0 || sc >= 1 {
		t.Fatalf("expected partial score, got %f", sc)
	}
}

func TestHybridRanks(t *testing.T) {
	s := NewScorer()
	memories := []Memory{
