package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Bookmark struct {
	ID           string     `json:"id"`
	CollectionID string     `json:"collection_id"`
	URL          string     `json:"url"`
	Title        string     `json:"title"`
	Tags         []string   `json:"tags"`
	ArchivedAt   *time.Time `json:"archived_at,omitempty"`
}

func NewBookmark(id, collectionID, rawURL, title string, tags []string) (Bookmark, error) {
	parsed, err := url.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || strings.TrimSpace(collectionID) == "" || strings.TrimSpace(title) == "" {
		return Bookmark{}, ErrInvalidBookmark
	}

	return Bookmark{
		ID:           id,
		CollectionID: strings.TrimSpace(collectionID),
		URL:          strings.TrimSpace(rawURL),
		Title:        strings.TrimSpace(title),
		Tags:         NormalizeTags(tags),
	}, nil
}

func BookmarkID(sequence int) string {
	return fmt.Sprintf("bookmark-%03d", sequence)
}

func NormalizeTags(tags []string) []string {
	return tags
}

func (b Bookmark) IsArchived() bool {
	return b.ArchivedAt != nil
}

func (b Bookmark) Clone() Bookmark {
	clone := b
	clone.Tags = append([]string(nil), b.Tags...)
	if b.ArchivedAt != nil {
		archivedAt := *b.ArchivedAt
		clone.ArchivedAt = &archivedAt
	}
	return clone
}
