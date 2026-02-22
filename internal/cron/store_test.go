package cron

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStore_SaveAndLoad(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	job := &Job{
		ID:        "test-1",
		Name:      "test job",
		Schedule:  "0 * * * * *",
		Tool:      "weather_current",
		Enabled:   true,
		CreatedAt: time.Now().Truncate(time.Second),
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}

	jobs, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].ID != "test-1" || jobs[0].Name != "test job" {
		t.Errorf("unexpected job: %+v", jobs[0])
	}
}

func TestStore_Delete(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	defer store.Close()

	job := &Job{
		ID:        "del-1",
		Name:      "to delete",
		Schedule:  "0 * * * * *",
		Enabled:   true,
		CreatedAt: time.Now(),
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatalf("SaveJob: %v", err)
	}
	if err := store.DeleteJob("del-1"); err != nil {
		t.Fatalf("DeleteJob: %v", err)
	}

	jobs, err := store.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs after delete, got %d", len(jobs))
	}
}
