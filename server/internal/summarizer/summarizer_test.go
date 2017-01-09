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
