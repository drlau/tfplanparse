package tfplanparse

import (
	"reflect"
	"testing"
)

func TestIsAttributeChangeLine(t *testing.T) {
	cases := map[string]struct {
		line     string
		expected bool
	}{
		"empty line": {
			line:     "",
			expected: false,
		},
		"attribute created": {
			line:     `+ attribute = "new"`,
			expected: true,
		},
		"attribute deleted": {
			line:     `- attribute = "deleted" -> null`,
			expected: true,
		},
		"attribute changed": {
			line:     `~ attribute = "old" -> "new"`,
			expected: true,
		},
		"attribute changed and forces replacement": {
			line:     `~ attribute = "old" -> "new" # forces replacement`,
			expected: true,
		},
		"attribute changed and value is unknown": {
			line:     `~ attribute = "old" -> (known after apply)`,
			expected: true,
		},
		"attribute is unchanged": {
			line:     `attribute = "old"`,
			expected: false,
		},
		"resource line": {
			line:     `+ resource "type" "name" {`,
			expected: false,
		},
		"padded with spaces": {
			line:     `    + attribute     = "new"`,
			expected: true,
		},
		"other line": {
			line:     `}`,
			expected: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := IsAttributeChangeLine(tc.line); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestNewAttributeChangeFromLine(t *testing.T) {
	cases := map[string]struct {
		line        string
		expected    *AttributeChange
		shouldError bool
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"attribute created": {
			line:        `+ attribute = "new"`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "",
				NewValue:   "new",
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- attribute = "deleted" -> null`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "deleted",
				NewValue:   "",
				UpdateType: DestroyResource,
			},
		},
		"attribute changed": {
			line:        `~ attribute = "old" -> "new"`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "old",
				NewValue:   "new",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute changed and forces replacement": {
			line:        `~ attribute = "old" -> "new" # forces replacement`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "old",
				NewValue:   "new",
				UpdateType: ForceReplaceResource,
			},
		},
		"attribute changed and value is unknown": {
			line:        `~ attribute = "old" -> (known after apply)`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "old",
				NewValue:   "(known after apply)",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute is unchanged": {
			line:        `attribute = "old"`,
			shouldError: true,
			expected:    nil,
		},
		"resource line": {
			line:        `+ resource "type" "name" {`,
			shouldError: true,
			expected:    nil,
		},
		"padded with spaces": {
			line:        `    + attribute     = "new"`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "",
				NewValue:   "new",
				UpdateType: NewResource,
			},
		},
		"other line": {
			line:        `}`,
			shouldError: true,
			expected:    nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := NewAttributeChangeFromLine(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
