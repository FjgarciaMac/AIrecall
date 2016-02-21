package storage

import (
	"os"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	path := os.TempDir() + "/airecall_store_test.db"
	_ = os.Remove(path)
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	if err := store.AddEpisode("a", "episode one"); err != nil {
		t.Fatal(err)
	}
