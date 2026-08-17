package service

import (
	"context"
	"sort"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/store"
)

type Library struct {
	repository store.Repository
	nextID     int
}

type BookmarkInput struct {
	CollectionID string   `json:"collection_id"`
	URL          string   `json:"url"`
	Title        string   `json:"title"`
	Tags         []string `json:"tags"`
}

type Report struct {
	ActiveByCollection map[string]int `json:"active_by_collection"`
	ArchivedCount      int            `json:"archived_count"`
}

func NewLibrary(repository store.Repository) *Library {
	return &Library{repository: repository}
}

func (l *Library) CreateCollection(ctx context.Context, name string) (domain.Collection, error) {
	l.nextID++
	collection, err := domain.NewCollection(domain.CollectionID(l.nextID), name)
	if err != nil {
		return domain.Collection{}, err
	}
	if err := l.repository.CreateCollection(ctx, collection); err != nil {
		return domain.Collection{}, err
	}
	return collection, nil
}

func (l *Library) ListCollections(ctx context.Context) ([]domain.Collection, error) {
	collections, err := l.repository.ListCollections(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(collections, func(i, j int) bool { return collections[i].ID < collections[j].ID })
	return collections, nil
}

func (l *Library) SaveBookmark(ctx context.Context, input BookmarkInput) (domain.Bookmark, error) {
	l.nextID++
	bookmark, err := domain.NewBookmark(domain.BookmarkID(l.nextID), input.CollectionID, input.URL, input.Title, input.Tags)
	if err != nil {
		return domain.Bookmark{}, err
	}
	if err := l.repository.SaveBookmark(ctx, bookmark); err != nil {
		return domain.Bookmark{}, err
	}
	return bookmark, nil
}

func (l *Library) MoveBookmark(ctx context.Context, bookmarkID, collectionID string) (domain.Bookmark, error) {
	return l.repository.MoveBookmark(ctx, bookmarkID, collectionID)
}

func (l *Library) ArchiveBookmark(ctx context.Context, bookmarkID string) (domain.Bookmark, error) {
	return l.repository.ArchiveBookmark(ctx, bookmarkID)
}

func (l *Library) BuildReport(ctx context.Context) (Report, error) {
	data, err := l.repository.BuildReport(ctx)
	if err != nil {
		return Report{}, err
	}
	report := Report{ActiveByCollection: make(map[string]int)}
	for _, collection := range data.Collections {
		report.ActiveByCollection[collection.ID] = 0
	}
	for _, bookmark := range data.Bookmarks {
		if bookmark.IsArchived() {
			report.ArchivedCount++
			continue
		}
		report.ActiveByCollection[bookmark.CollectionID]++
	}
	return report, nil
}
