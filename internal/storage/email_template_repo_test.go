//go:build integration

package storage

import (
	"testing"
)

func TestEmailTemplateRepo_SeedUpsertGet(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewEmailTemplateRepository(pool)

	if err := repo.SeedIfMissing("t_seed", "Seed", "Subj", "<p>hi</p>", "hi", "desc", []string{"NAME"}); err != nil {
		t.Fatalf("SeedIfMissing: %v", err)
	}
	// seed does not overwrite
	if err := repo.SeedIfMissing("t_seed", "Other", "X", "x", "x", "x", nil); err != nil {
		t.Fatalf("SeedIfMissing 2: %v", err)
	}
	got, err := repo.GetByKey("t_seed")
	if err != nil {
		t.Fatalf("GetByKey: %v", err)
	}
	if got.Name != "Seed" || got.IsModified {
		t.Errorf("seed overwrote: %+v", got)
	}

	if err := repo.Upsert("t_seed", "Edited", "Subj2", "<p>bye</p>", "bye", "d2", []string{"NAME", "URL"}); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	got, _ = repo.GetByKey("t_seed")
	if got.Name != "Edited" || !got.IsModified || len(got.RequiredVars) != 2 {
		t.Errorf("upsert: %+v", got)
	}

	if err := repo.Delete("t_seed"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.GetByKey("t_seed"); err == nil {
		t.Error("expected not found after delete")
	}
}

func TestEmailTemplateRepo_List(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewEmailTemplateRepository(pool)

	_ = repo.Upsert("t_list_a", "A", "s", "h", "t", "d", nil)
	_ = repo.Upsert("t_list_b", "B", "s", "h", "t", "d", nil)
	list, err := repo.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	found := map[string]bool{}
	for _, e := range list {
		found[e.Key] = true
	}
	if !found["t_list_a"] || !found["t_list_b"] {
		t.Errorf("list missing keys: %v", found)
	}
}
