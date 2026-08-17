package domain

import (
	"fmt"
	"strings"
)

type Collection struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewCollection(id, name string) (Collection, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(name) == "" {
		return Collection{}, ErrInvalidCollection
	}

	return Collection{ID: strings.TrimSpace(id), Name: strings.TrimSpace(name)}, nil
}

func CollectionID(sequence int) string {
	return fmt.Sprintf("collection-%03d", sequence)
}
