package summarizer

import (
	"os"
	"testing"

	"github.com/FjgarciaMac/AIrecall/server/internal/storage"
)

func TestSplitFact(t *testing.T) {
	cases := map[string][]string{
		"prefers email":        {"prefers", "email"},
		"contact: email":       {"contact", "email"},
		"timezone is UTC":      {"timezone", "UTC"},
		"name = alice":         {"name", "alice"},
		"a very long key here": nil,
	}
	for in, want := range cases {
		got := splitFact(in)
		if want == nil {
			if got != nil {
				t.Fatalf("%q: expected nil, got %v", in, got)
			}
			continue
		}
		if got == nil || got[0] != want[0] || got[1] != want[1] {
			t.Fatalf("%q: got %v want %v", in, got, want)
		}
	}
}

func TestRunPromotesFacts(t *testing.T) {
	path := os.TempDir() + "/airecall_summarizer_test.db"
	_ = os.Remove(path)
	store, err := storage.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	_ = store.AddEpisode("a", "timezone is UTC")
	_ = store.AddEpisode("a", "no fact in this line at all")
	res, err := New(store).Run("a", 10)
	if err != nil {
		t.Fatal(err)
	}
