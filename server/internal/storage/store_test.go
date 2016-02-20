package storage

import (
	"os"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	path := os.TempDir() + "/airecall_store_test.db"
