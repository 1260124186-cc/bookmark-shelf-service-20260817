package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/service"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
)

func TestLibraryOrganizesAndReportsBookmarks(t *testing.T) {
	t.Parallel()
	library := service.NewLibrary(store.NewMemoryRepository())
	ctx := context.Background()

	work, err := library.CreateCollection(ctx, "Work")
	if err != nil {
		t.Fatal(err)
	}
	reading, err := library.CreateCollection(ctx, "Reading")
	if err != nil {
		t.Fatal(err)
	}
	bookmark, err := library.SaveBookmark(ctx, service.BookmarkInput{
		CollectionID: work.ID,
		URL:          "https://go.dev/doc/",
		Title:        "Go documentation",
		Tags:         []string{"Go", " reference ", "go"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := library.MoveBookmark(ctx, bookmark.ID, reading.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := library.ArchiveBookmark(ctx, bookmark.ID); err != nil {
		t.Fatal(err)
	}

	report, err := library.BuildReport(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if report.ArchivedCount != 1 || report.ActiveByCollection[reading.ID] != 0 {
		t.Fatalf("unexpected report: %#v", report)
	}
}

func TestLibraryRejectsDuplicateWithinCollection(t *testing.T) {
	t.Parallel()
	library := service.NewLibrary(store.NewMemoryRepository())
	ctx := context.Background()
	collection, err := library.CreateCollection(ctx, "Research")
	if err != nil {
		t.Fatal(err)
	}
	input := service.BookmarkInput{CollectionID: collection.ID, URL: "https://example.com/article", Title: "Article"}
	if _, err := library.SaveBookmark(ctx, input); err != nil {
		t.Fatal(err)
	}
	if _, err := library.SaveBookmark(ctx, input); !errors.Is(err, domain.ErrDuplicateBookmark) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestLibraryKeepsNormalizedTagsAfterCallerReusesInput(t *testing.T) {
	t.Parallel()
	library := service.NewLibrary(store.NewMemoryRepository())
	ctx := context.Background()
	collection, err := library.CreateCollection(ctx, "Research")
	if err != nil {
		t.Fatal(err)
	}
	tags := []string{" Go ", "reference", "go"}
	bookmark, err := library.SaveBookmark(ctx, service.BookmarkInput{
		CollectionID: collection.ID,
		URL:          "https://example.com/go",
		Title:        "Go",
		Tags:         tags,
	})
	if err != nil {
		t.Fatal(err)
	}
	tags[0] = "changed"
	if bookmark.Tags[0] != "go" || len(bookmark.Tags) != 2 {
		t.Fatalf("unexpected returned tags: %#v", bookmark.Tags)
	}
}
