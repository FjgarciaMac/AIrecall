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
