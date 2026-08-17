package store

import (
	"context"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
)

type Repository interface {
	CreateCollection(context.Context, domain.Collection) error
	ListCollections(context.Context) ([]domain.Collection, error)
	SaveBookmark(context.Context, domain.Bookmark) error
	GetBookmark(context.Context, string) (domain.Bookmark, error)
	MoveBookmark(context.Context, string, string) (domain.Bookmark, error)
	ArchiveBookmark(context.Context, string) (domain.Bookmark, error)
	BuildReport(context.Context) (ReportData, error)
}

type ReportData struct {
	Collections []domain.Collection
	Bookmarks   []domain.Bookmark
}
