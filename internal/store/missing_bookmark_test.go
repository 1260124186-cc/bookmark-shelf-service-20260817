package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
)

func TestMemoryRepositoryReturnsMissingBookmarkError(t *testing.T) {
	t.Parallel()
	repository := store.NewMemoryRepository()
	_, err := repository.GetBookmark(context.Background(), "bookmark-missing")
	if !errors.Is(err, domain.ErrBookmarkNotFound) {
		t.Fatalf("expected missing bookmark error, got %v", err)
	}
}
