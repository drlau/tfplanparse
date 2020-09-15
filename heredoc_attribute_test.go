package tfplanparse

import (
	"reflect"
	"testing"
)

func TestNewHeredocAttributeChangeFromLine(t *testing.T) {
	cases := map[string]struct {
		line        string
		expected    *HeredocAttributeChange
		shouldError bool
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"attribute created": {
			line:        `+ attribute = <<~EOT`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- attribute = <<~EOT`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: DestroyResource,
			},
		},
		"attribute changed": {
			line:        `~ attribute = <<~EOT`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute changed and forced replacement": {
			line:        `~ attribute = <<~EOT # forces replacement`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: ForceReplaceResource,
			},
		},
		"attribute is unchanged": {
			line:        `attribute = <<~EOT`,
			shouldError: true,
			expected:    nil,
		},
		"malformed heredoc": {
			line:        `+ attribute = <EOT`,
			shouldError: true,
			expected:    nil,
		},
		"tilda heredoc": {
			line:        `+ attribute = <<~EOT`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: NewResource,
			},
		},
		"resource line": {
			line:        `+ resource "type" "name" {`,
			shouldError: true,
			expected:    nil,
		},
		"padded with spaces": {
			line:        `    + attribute     = <<~EOT`,
			shouldError: false,
			expected: &HeredocAttributeChange{
				Name:       "attribute",
				Before:     []string{},
				After:      []string{},
				UpdateType: NewResource,
			},
		},
		"other line": {
			line:        `]`,
			shouldError: true,
			expected:    nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := NewHeredocAttributeChangeFromLine(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
