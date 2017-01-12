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
