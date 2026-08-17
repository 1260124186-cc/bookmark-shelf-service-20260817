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

func TestLibraryStopsReportWhenRequestIsCanceled(t *testing.T) {
	t.Parallel()
	library := service.NewLibrary(store.NewMemoryRepository())
	ctx := context.Background()
	if _, err := library.CreateCollection(ctx, "Reading"); err != nil {
		t.Fatal(err)
	}

	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := library.BuildReport(canceled); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled report, got %v", err)
	}
}
