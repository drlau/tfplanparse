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
				OldValue:   nil,
				NewValue:   "new",
				UpdateType: NewResource,
			},
		},
		"bools are parsed as bools": {
			line:        `+ attribute = true`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   nil,
				NewValue:   true,
				UpdateType: NewResource,
			},
		},
		"ints are parsed as ints": {
			line:        `+ attribute = 1`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   nil,
				NewValue:   1,
				UpdateType: NewResource,
			},
		},
		"decimals are parsed as floats": {
			line:        `+ attribute = 1.23`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   nil,
				NewValue:   1.23,
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- attribute = "deleted" -> null`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "deleted",
				NewValue:   nil,
				UpdateType: DestroyResource,
			},
		},
		"empty map is deleted": {
			line:        `- attribute {}`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   nil,
				NewValue:   nil,
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
		"attribute changed and value is sensitive": {
			line:        `~ attribute = (sensitive value)`,
			shouldError: false,
			expected: &AttributeChange{
				Name:       "attribute",
				OldValue:   "(sensitive value)",
				NewValue:   "(sensitive value)",
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
				OldValue:   nil,
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

func TestNewAttributeChangeFromArray(t *testing.T) {
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
			line:        `+ "new"`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
				NewValue:   "new",
				UpdateType: NewResource,
			},
		},
		"bools are parsed as bools": {
			line:        `+ true`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
				NewValue:   true,
				UpdateType: NewResource,
			},
		},
		"ints are parsed as ints": {
			line:        `+ 1`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
				NewValue:   1,
				UpdateType: NewResource,
			},
		},
		"decimals are parsed as floats": {
			line:        `+ 1.23`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
				NewValue:   1.23,
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- "deleted"`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   "deleted",
				NewValue:   nil,
				UpdateType: DestroyResource,
			},
		},
		"empty map is deleted": {
			line:        `- {}`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
				NewValue:   nil,
				UpdateType: DestroyResource,
			},
		},
		"attribute is unchanged": {
			line:        `"old"`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   "old",
				NewValue:   "old",
				UpdateType: NoOpResource,
			},
		},
		"resource line": {
			line:        `+ resource "type" "name" {`,
			shouldError: true,
			expected:    nil,
		},
		"padded with spaces": {
			line:        `    +        "new"`,
			shouldError: false,
			expected: &AttributeChange{
				OldValue:   nil,
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
			got, err := NewAttributeChangeFromArray(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
