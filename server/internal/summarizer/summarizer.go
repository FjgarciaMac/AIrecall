// Package summarizer distills durable facts from older episodes.
package summarizer

import (
	"regexp"
	"strings"

	"github.com/FjgarciaMac/AIrecall/server/internal/storage"
)

// Extractor finds "key: value" style facts in episode text.
var factPattern = regexp.MustCompile(
