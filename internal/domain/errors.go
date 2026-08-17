package domain

import "errors"

var (
	ErrInvalidCollection  = errors.New("invalid collection")
	ErrInvalidBookmark    = errors.New("invalid bookmark")
	ErrCollectionNotFound = errors.New("collection not found")
	ErrBookmarkNotFound   = errors.New("bookmark not found")
	ErrDuplicateBookmark  = errors.New("bookmark already exists in collection")
)
