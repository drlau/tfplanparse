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
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NoOpResource,
			},
		},
		"attribute with delimiter is unchanged": {
			line:        `attribute = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NoOpResource,
			},
		},
		"unchanged empty map": {
			line:        `attribute = {}`,
			shouldError: false,
			expected: &MapAttributeChange{
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
			line:        `    + attribute     = {`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"one line empty map": {
			line:        `+ attribute = {}`,
			shouldError: false,
			expected: &MapAttributeChange{
				Name:       "attribute",
				UpdateType: NewResource,
			},
		},
		"one line empty map with no delimiter": {
			line:        `+ attribute {}`,
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

func TestMapGetBefore(t *testing.T) {
	cases := map[string]struct {
		ma       *MapAttributeChange
		expected map[string]interface{}
		opts     []GetBeforeAfterOptions
	}{
		"one attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
				},
			},
			expected: map[string]interface{}{
				"attribute": "oldValue",
			},
		},
		"multiple attribute": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"attribute1": "oldValue1",
				"attribute2": "oldValue2",
			},
		},
		"map attribute": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"attribute1": "oldValue1",
					"attribute2": "oldValue2",
				},
			},
		},
		"array attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&ArrayAttributeChange{
						Name: "array",
						AttributeChanges: []attributeChange{
							&AttributeChange{
								OldValue: "oldValue1",
								NewValue: "newValue1",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"array": []interface{}{
					"oldValue1",
				},
			},
		},
		"map and normal attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
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
			expected: map[string]interface{}{
				"attribute": "oldValue",
				"map": map[string]interface{}{
					"attribute1": "oldValue1",
					"attribute2": "oldValue2",
				},
			},
		},
		"ignore sensitive values": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"attribute2": "oldValue2",
			},
			opts: []GetBeforeAfterOptions{IgnoreSensitive},
		},
		// no tests for computed options because "before" values are never "computed"
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.ma.GetBefore(tc.opts...); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestMapGetAfter(t *testing.T) {
	cases := map[string]struct {
		ma       *MapAttributeChange
		expected map[string]interface{}
		opts     []GetBeforeAfterOptions
	}{
		"one attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
				},
			},
			expected: map[string]interface{}{
				"attribute": "newValue",
			},
		},
		"multiple attribute": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"attribute1": "newValue1",
				"attribute2": "newValue2",
			},
		},
		"map attribute": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"map": map[string]interface{}{
					"attribute1": "newValue1",
					"attribute2": "newValue2",
				},
			},
		},
		"array attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&ArrayAttributeChange{
						Name: "array",
						AttributeChanges: []attributeChange{
							&AttributeChange{
								OldValue: "oldValue1",
								NewValue: "newValue1",
							},
						},
					},
				},
			},
			expected: map[string]interface{}{
				"array": []interface{}{
					"newValue1",
				},
			},
		},
		"map and normal attribute": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
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
			expected: map[string]interface{}{
				"attribute": "newValue",
				"map": map[string]interface{}{
					"attribute1": "newValue1",
					"attribute2": "newValue2",
				},
			},
		},
		"ignore sensitive values": {
			ma: &MapAttributeChange{
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
			expected: map[string]interface{}{
				"attribute2": "newValue2",
			},
			opts: []GetBeforeAfterOptions{IgnoreSensitive},
		},
		"ignore computed values": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "(known after apply)",
					},
					&AttributeChange{
						Name:     "attribute2",
						OldValue: "oldValue2",
						NewValue: "newValue2",
					},
				},
			},
			expected: map[string]interface{}{
				"attribute2": "newValue2",
			},
			opts: []GetBeforeAfterOptions{IgnoreComputed},
		},
		"computed only": {
			ma: &MapAttributeChange{
				AttributeChanges: []attributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "(known after apply)",
					},
					&AttributeChange{
						Name:     "attribute2",
						OldValue: "oldValue2",
						NewValue: "newValue2",
					},
				},
			},
			expected: map[string]interface{}{
				"attribute": "(known after apply)",
			},
			opts: []GetBeforeAfterOptions{ComputedOnly},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := tc.ma.GetAfter(tc.opts...); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
