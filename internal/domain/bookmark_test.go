package domain

import (
	"reflect"
	"testing"
)

func TestNormalizeTagsCreatesIndependentCanonicalTags(t *testing.T) {
	t.Parallel()

	tags := []string{" Go ", "reference", "go", " ", "REFERENCE", "Testing "}
	normalized := NormalizeTags(tags)

	if want := []string{"go", "reference", "testing"}; !reflect.DeepEqual(normalized, want) {
		t.Fatalf("NormalizeTags() = %#v, want %#v", normalized, want)
	}

	tags[0] = "changed"
	if normalized[0] != "go" {
		t.Fatalf("normalized tags changed with caller input: %#v", normalized)
	}
}
