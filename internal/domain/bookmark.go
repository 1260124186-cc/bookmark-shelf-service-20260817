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

// NormalizeTags 规范化标签：去空白、转小写、去重，并返回独立切片，
// 避免调用方复用或修改原切片时影响已保存的书签数据。
func NormalizeTags(tags []string) []string {
	seen := make(map[string]struct{}, len(tags))
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		cleaned := strings.ToLower(strings.TrimSpace(tag))
		if cleaned == "" {
			continue
		}
		if _, exists := seen[cleaned]; exists {
			continue
		}
		seen[cleaned] = struct{}{}
		normalized = append(normalized, cleaned)
	}
	return normalized
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
