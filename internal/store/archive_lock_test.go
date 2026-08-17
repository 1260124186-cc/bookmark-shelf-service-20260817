package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
)

func TestArchiveMissDoesNotBlockLaterWrites(t *testing.T) {
	repository := store.NewMemoryRepository()
	if _, err := repository.ArchiveBookmark(context.Background(), "bookmark-missing"); err == nil {
		t.Fatal("expected missing bookmark error")
	}
	collection, err := domain.NewCollection("collection-001", "Reading")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- repository.CreateCollection(context.Background(), collection) }()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("later write remained blocked after archive miss")
	}
}
