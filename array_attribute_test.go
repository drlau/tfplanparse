package tfplanparse

import (
	"reflect"
	"testing"
)

func TestNewArrayAttributeChangeFromLine(t *testing.T) {
	cases := map[string]struct {
		line        string
		expected    *ArrayAttributeChange
		shouldError bool
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"attribute created": {
			line:        `+ attribute [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"attribute created with delimiter": {
			line:        `+ attribute = [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"attribute deleted": {
			line:        `- attribute [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: DestroyResource,
			},
		},
		"attribute deleted with delimiter": {
			line:        `- attribute = [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: DestroyResource,
			},
		},
		"attribute changed": {
			line:        `~ attribute [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute changed with delimiter": {
			line:        `~ attribute = [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"attribute is unchanged": {
			line:        `attribute [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NoOpResource,
			},
		},
		"attribute with delimiter is unchanged": {
			line:        `attribute = [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NoOpResource,
			},
		},
		"unchanged empty array": {
			line:        `attribute = []`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NoOpResource,
			},
		},
		"resource line": {
			line:        `+ resource "type" "name" {`,
			shouldError: true,
			expected:    nil,
		},
		"padded with spaces": {
			line:        `    + attribute     = [`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"one line empty array": {
			line:        `+ attribute = []`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"one line empty array with no delimiter": {
			line:        `+ attribute []`,
			shouldError: false,
			expected: &ArrayAttributeChange{
				Name:       "attribute",
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
			got, err := NewArrayAttributeChangeFromLine(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestGetBefore(t *testing.T) {
	cases := map[string]struct {
		aa       *ArrayAttributeChange
		expected []interface{}
		opts     []GetBeforeAfterOptions
	}{
		"one attribute": {
			aa: &ArrayAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
				},
			},
			expected: []interface{}{
				"oldValue",
			},
		},
		"multiple attribute": {
			aa: &ArrayAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute1",
						OldValue: "oldValue1",
						NewValue: "newValue1",
					},
					&AttributeChange{
						Name:     "attribute2",
						OldValue: "oldValue2",
						NewValue: "newValue2",
					},
				},
			},
			expected: []interface{}{
				"oldValue1",
				"oldValue2",
			},
		},
		"ignores new attributes": {
			aa: &ArrayAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute1",
						OldValue: "oldValue1",
						NewValue: "newValue1",
					},
					&AttributeChange{
						Name:       "attribute2",
						OldValue:   nil,
						NewValue:   "newValue2",
						UpdateType: NewResource,
					},
				},
			},
			expected: []interface{}{
				"oldValue1",
			},
		},
		"map attribute": {
			aa: &ArrayAttributeChange{
				AttributeChanges: []attributeChange{
					&MapAttributeChange{
						Name: "map",
						AttributeChanges: []attributeChange{
							&AttributeChange{
								Name:     "attribute1",
								OldValue: "oldValue1",
								NewValue: "newValue1",
							},
							&AttributeChange{
								Name:     "attribute2",
								OldValue: "oldValue2",
								NewValue: "newValue2",
							},
						},
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"attribute1": "oldValue1",
					"attribute2": "oldValue2",
				},
			},
		},
		"ignore sensitive values": {
			aa: &ArrayAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "(sensitive value)",
						NewValue: "(sensitive value)",
					},
					&AttributeChange{
						Name:     "attribute2",
						OldValue: "oldValue2",
						NewValue: "newValue2",
					},
				},
			},
			expected: []interface{}{
				"oldValue2",
			},
			opts: []GetBeforeAfterOptions{IgnoreSensitive},
		},
		// no tests for computed options because "before" values are never "computed"
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.aa.GetBefore(tc.opts...); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
