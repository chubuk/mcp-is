//go:build integration

package storage

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
)

func TestIdempotencyRepoReplayAndReuse(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewIdempotencyRepository(pool)
	ctx := context.Background()
	keyID := "test-idem-key-" + t.Name()
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM idempotency_keys WHERE ai_key_id=$1`, keyID)
	})

	// Begin → created=true (pending).
	rec, created, err := repo.Begin(ctx, keyID, "k1", "create_product", "hash1")
	if err != nil || !created {
		t.Fatalf("first Begin: rec=%v created=%v err=%v", rec, created, err)
	}

	// Параллельный повтор до Finish → pending, не created.
	rec2, created2, err := repo.Begin(ctx, keyID, "k1", "create_product", "hash1")
	if err != nil || created2 {
		t.Fatalf("concurrent Begin: created=%v err=%v", created2, err)
	}
	if rec2.Status != "pending" {
		t.Errorf("expected pending, got %q", rec2.Status)
	}

	// Finish → replay c done.
	payload := []byte(`{"text":"{\"id\":1}","structured":{"id":1}}`)
	if err := repo.Finish(ctx, keyID, "k1", payload); err != nil {
		t.Fatal(err)
	}
	rec3, created3, err := repo.Begin(ctx, keyID, "k1", "create_product", "hash1")
	if err != nil || created3 {
		t.Fatalf("replay Begin: created=%v err=%v", created3, err)
	}
	if rec3.Status != "done" {
		t.Errorf("expected done, got %q", rec3.Status)
	}
	// JSONB нормализует пробелы — сравниваем семантически.
	var got, want map[string]any
	if jerr := json.Unmarshal(rec3.Response, &got); jerr != nil {
		t.Fatalf("replay response not JSON: %v", jerr)
	}
	if jerr := json.Unmarshal(payload, &want); jerr != nil {
		t.Fatal(jerr)
	}
	if got["text"] != want["text"] {
		t.Errorf("replay payload mismatch: %v != %v", got, want)
	}

	// Тот же ключ с другим hash → ошибка переиспользования.
	_, created4, err := repo.Begin(ctx, keyID, "k1", "create_product", "hash2")
	if err != nil || created4 {
		t.Fatalf("reuse Begin: created=%v err=%v", created4, err)
	}
}

func TestIdempotencyRepoCancelAllowsRetry(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewIdempotencyRepository(pool)
	ctx := context.Background()
	keyID := "test-idem-cancel-" + t.Name()
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM idempotency_keys WHERE ai_key_id=$1`, keyID)
	})

	if _, created, err := repo.Begin(ctx, keyID, "k1", "create_order", "h"); err != nil || !created {
		t.Fatalf("Begin: %v %v", created, err)
	}
	if err := repo.Cancel(ctx, keyID, "k1"); err != nil {
		t.Fatal(err)
	}
	_, created, err := repo.Begin(ctx, keyID, "k1", "create_order", "h")
	if err != nil || !created {
		t.Fatalf("retry after cancel must create: created=%v err=%v", created, err)
	}
}

func TestIdempotencyRepoKeyScoping(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewIdempotencyRepository(pool)
	ctx := context.Background()
	a, b := "test-idem-scope-a", "test-idem-scope-b"
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM idempotency_keys WHERE ai_key_id IN ($1,$2)`, a, b)
	})

	if _, created, err := repo.Begin(ctx, a, "k", "create_product", "h"); err != nil || !created {
		t.Fatalf("Begin a: %v %v", created, err)
	}
	// Один и тот же ключ у другого AI-ключа — независимая запись.
	if _, created, err := repo.Begin(ctx, b, "k", "create_product", "h"); err != nil || !created {
		t.Fatalf("Begin b must be independent: created=%v err=%v", created, err)
	}
}

// TestIdempotencyRepoConcurrentBegin — гонка параллельных retry одного ключа:
// INSERT ... ON CONFLICT (PK) гарантирует ровно один created=true, остальные
// видят pending и получают «in progress» на уровне обёртки (не дубль).
func TestIdempotencyRepoConcurrentBegin(t *testing.T) {
	pool := setupTestDB(t)
	repo := NewIdempotencyRepository(pool)
	ctx := context.Background()
	keyID := "test-idem-race"
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM idempotency_keys WHERE ai_key_id = $1`, keyID)
	})

	const n = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]bool, n) // created per goroutine
	errs := make([]error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_, created, err := repo.Begin(ctx, keyID, "k-race", "create_order", "hash")
			results[i] = created
			errs[i] = err
		}(i)
	}
	close(start)
	wg.Wait()

	winners := 0
	for i := 0; i < n; i++ {
		if errs[i] != nil {
			t.Fatalf("goroutine %d: %v", i, errs[i])
		}
		if results[i] {
			winners++
		}
	}
	if winners != 1 {
		t.Fatalf("created=true count = %d, want exactly 1", winners)
	}
}
