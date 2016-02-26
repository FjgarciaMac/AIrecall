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
	eps, err := store.RecentEpisodes("a", 10)
	if err != nil || len(eps) != 1 || eps[0].Content != "episode one" {
		t.Fatalf("episodes = %v err=%v", eps, err)
	}

	if err := store.UpsertFact("a", "prefers", "email"); err != nil {
		t.Fatal(err)
	}
	if err := store.UpsertFact("a", "prefers", "phone"); err != nil {
		t.Fatal(err)
	}
	value, ok, err := store.Fact("a", "prefers")
	if err != nil || !ok || value != "phone" {
