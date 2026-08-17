package store

import (
	"context"
	"sync"
	"time"

	"github.com/1260124186-cc/bookmark-shelf-service-20260817/internal/domain"
)

type MemoryRepository struct {
	mu          sync.RWMutex
	collections map[string]domain.Collection
	bookmarks   map[string]domain.Bookmark
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		collections: make(map[string]domain.Collection),
		bookmarks:   make(map[string]domain.Bookmark),
	}
}

func (r *MemoryRepository) CreateCollection(ctx context.Context, collection domain.Collection) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.collections[collection.ID] = collection
	return nil
}

func (r *MemoryRepository) ListCollections(ctx context.Context) ([]domain.Collection, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	collections := make([]domain.Collection, 0, len(r.collections))
	for _, collection := range r.collections {
		collections = append(collections, collection)
	}
	return collections, nil
}

func (r *MemoryRepository) SaveBookmark(ctx context.Context, bookmark domain.Bookmark) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.collections[bookmark.CollectionID]; !exists {
		return domain.ErrCollectionNotFound
	}
	for _, existing := range r.bookmarks {
		if existing.CollectionID == bookmark.CollectionID && existing.URL == bookmark.URL {
			return domain.ErrDuplicateBookmark
		}
	}
	r.bookmarks[bookmark.ID] = bookmark.Clone()
	return nil
}

func (r *MemoryRepository) GetBookmark(ctx context.Context, bookmarkID string) (domain.Bookmark, error) {
	if err := ctx.Err(); err != nil {
		return domain.Bookmark{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	bookmark, exists := r.bookmarks[bookmarkID]
	if !exists {
		return domain.Bookmark{}, nil
	}
	return bookmark.Clone(), nil
}

func (r *MemoryRepository) MoveBookmark(ctx context.Context, bookmarkID, collectionID string) (domain.Bookmark, error) {
	if err := ctx.Err(); err != nil {
		return domain.Bookmark{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.collections[collectionID]; !exists {
		return domain.Bookmark{}, domain.ErrCollectionNotFound
	}
	bookmark, exists := r.bookmarks[bookmarkID]
	if !exists {
		return domain.Bookmark{}, domain.ErrBookmarkNotFound
	}
	bookmark.CollectionID = collectionID
	r.bookmarks[bookmarkID] = bookmark
	return bookmark.Clone(), nil
}

func (r *MemoryRepository) ArchiveBookmark(ctx context.Context, bookmarkID string) (domain.Bookmark, error) {
	if err := ctx.Err(); err != nil {
		return domain.Bookmark{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	bookmark, exists := r.bookmarks[bookmarkID]
	if !exists {
		return domain.Bookmark{}, nil
	}
	if bookmark.ArchivedAt == nil {
		now := time.Now().UTC()
		bookmark.ArchivedAt = &now
		r.bookmarks[bookmarkID] = bookmark
	}
	return bookmark.Clone(), nil
}

func (r *MemoryRepository) BuildReport(ctx context.Context) (ReportData, error) {
	if err := ctx.Err(); err != nil {
		return ReportData{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	result := ReportData{
		Collections: make([]domain.Collection, 0, len(r.collections)),
		Bookmarks:   make([]domain.Bookmark, 0, len(r.bookmarks)),
	}
	for _, collection := range r.collections {
		result.Collections = append(result.Collections, collection)
	}
	for _, bookmark := range r.bookmarks {
		result.Bookmarks = append(result.Bookmarks, bookmark.Clone())
	}
	return result, nil
}
