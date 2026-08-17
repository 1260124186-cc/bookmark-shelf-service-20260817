package store_test

import (
	"context"
	"testing"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
)

func TestMemoryRepositoryDoesNotRetainCallerTagSlice(t *testing.T) {
	t.Parallel()
	repository := store.NewMemoryRepository()
	ctx := context.Background()
	collection, err := domain.NewCollection("collection-001", "Work")
	if err != nil {
		t.Fatal(err)
	}
	if err := repository.CreateCollection(ctx, collection); err != nil {
		t.Fatal(err)
	}
	bookmark := domain.Bookmark{
		ID:           "bookmark-001",
		CollectionID: collection.ID,
		URL:          "https://example.com",
		Title:        "Example",
		Tags:         []string{"work"},
	}
	if err := repository.SaveBookmark(ctx, bookmark); err != nil {
		t.Fatal(err)
	}
	bookmark.Tags[0] = "changed"
	saved, err := repository.GetBookmark(ctx, bookmark.ID)
	if err != nil {
		t.Fatal(err)
	}
	if saved.Tags[0] != "work" {
		t.Fatalf("stored tags changed with caller input: %#v", saved.Tags)
	}
}
