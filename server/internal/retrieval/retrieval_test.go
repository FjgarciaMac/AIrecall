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
		{Content: "unrelated note about weather"},
		{Content: "user prefers email over phone"},
		{Content: "told user about refund policy yesterday"},
		{Content: "refund policy applies within 30 days"},
	}
	hits := s.Hybrid("refund policy", memories, 2, "hybrid")
	if len(hits) != 2 {
		t.Fatalf("expected 2 hits, got %d", len(hits))
	}
	if hits[0] != memories[3].Content {
		t.Fatalf("expected most relevant first, got %q", hits[0])
	}
	if hits[1] != memories[2].Content {
		t.Fatalf("expected second relevant, got %q", hits[1])
	}
}

func TestVectorMode(t *testing.T) {
	s := NewScorer()
	sc := s.Score("email contact", "prefers email contact", "vector")
	if sc <= 0 {
		t.Fatalf("expected vector similarity > 0, got %f", sc)
	}
}

func TestStopwordsIgnored(t *testing.T) {
	s := NewScorer()
	sc := s.Score("a b the", "nothing matches here", "keyword")
	if sc != 0 {
		t.Fatalf("stopwords should score 0, got %f", sc)
	}
}
// draft note 1402
