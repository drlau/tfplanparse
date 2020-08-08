package tfplanparse

import (
	"reflect"
	"testing"
)

func TestNewMapAttributeChangeFromLine(t *testing.T) {
	cases := map[string]struct {
		line        string
		expected    *MapAttributeChange
		shouldError bool
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"attribute created": {
			line:        `+ attribute {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"attribute created with delimiter": {
			line:        `+ attribute = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- attribute {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: DestroyResource,
			},
		},
		"attribute deleted with delimiter": {
			line:        `- attribute = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: DestroyResource,
			},
		},
		"attribute changed": {
			line:        `~ attribute {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute changed with delimiter": {
			line:        `~ attribute = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute is unchanged": {
			line:        `attribute {`,
			shouldError: true,
			expected:    nil,
		},
		"attribute with delimiter is unchanged": {
			line:        `attribute = {`,
			shouldError: true,
			expected:    nil,
		},
		"resource line": {
			line:        `+ resource "type" "name" {`,
			shouldError: true,
			expected:    nil,
		},
		"padded with spaces": {
			line:        `    + attribute     = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
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
			got, err := NewMapAttributeChangeFromLine(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
