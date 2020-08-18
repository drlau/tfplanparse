package tfplanparse

import (
	"reflect"
	"testing"
)

func TestIsResourceCommentLine(t *testing.T) {
	cases := map[string]struct {
		line     string
		expected bool
	}{
		"empty line": {
			line:     "",
			expected: false,
		},
		"resource created": {
			line:     "# resource.path will be created",
			expected: true,
		},
		"resource read during apply": {
			line:     "# resource.path will be created",
			expected: true,
		},
		"resource read during apply extra line": {
			line:     "# (config refers to values not yet known)",
			expected: false,
		},
		"resource updated in place": {
			line:     "# resource.path will be updated in-place",
			expected: true,
		},
		"resource tainted": {
			line:     "# resource.path is tainted, so must be replaced",
			expected: true,
		},
		"resource replaced": {
			line:     "# resource.path must be replaced",
			expected: true,
		},
		"resource destroyed": {
			line:     "# resource.path will be destroyed",
			expected: true,
		},
		"handles extra spaces": {
			line:     "    # resource.path will be created",
			expected: true,
		},
		"other line": {
			line:     "~ resource",
			expected: false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := IsResourceCommentLine(tc.line); got != tc.expected {
				t.Errorf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestNewResourceChangeFromComment(t *testing.T) {
	cases := map[string]struct {
		line        string
		shouldError bool
		expected    *ResourceChange
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"resource created": {
			line:        "# resource.path will be created",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: NewResource,
			},
		},
		"resource read during apply": {
			line:        "# resource.path will be read during apply",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ReadResource,
			},
		},
		"resource read during apply extra line": {
			line:        "# (config refers to values not yet known)",
			shouldError: true,
			expected:    nil,
		},
		"resource updated in place": {
			line:        "# resource.path will be updated in-place",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"resource tainted": {
			line:        "# resource.path is tainted, so must be replaced",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ForceReplaceResource,
				Tainted:    true,
			},
		},
		"resource replaced": {
			line:        "# resource.path must be replaced",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ForceReplaceResource,
			},
		},
		"resource destroyed": {
			line:        "# resource.path will be destroyed",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: DestroyResource,
			},
		},
		"handles extra spaces": {
			line:        "    # resource.path will be created",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: NewResource,
			},
		},
		"handles data": {
			line:        "    # data.mydata.path will be read during apply",
			shouldError: false,
			expected: &ResourceChange{
				Address:    "data.mydata.path",
				Type:       "mydata",
				Name:       "path",
				UpdateType: ReadResource,
			},
		},
		"handles modules with extra spaces": {
			line:        "    # module.mymodule.resource.path will be created",
			shouldError: false,
			expected: &ResourceChange{
				Address:       "module.mymodule.resource.path",
				ModuleAddress: "module.mymodule",
				Type:          "resource",
				Name:          "path",
				UpdateType:    NewResource,
			},
		},
		"string index": {
			line:        `    # module.mymodule.resource.path["index"] will be created`,
			shouldError: false,
			expected: &ResourceChange{
				Address:       `module.mymodule.resource.path["index"]`,
				ModuleAddress: "module.mymodule",
				Type:          "resource",
				Name:          "path",
				Index:         "index",
				UpdateType:    NewResource,
			},
		},
		"string index with a .": {
			line:        `    # module.mymodule.resource.path["index@test.com"] will be created`,
			shouldError: false,
			expected: &ResourceChange{
				Address:       `module.mymodule.resource.path["index@test.com"]`,
				ModuleAddress: "module.mymodule",
				Type:          "resource",
				Name:          "path",
				Index:         "index@test.com",
				UpdateType:    NewResource,
			},
		},
		"int index": {
			line:        "    # module.mymodule.resource.path[1] will be created",
			shouldError: false,
			expected: &ResourceChange{
				Address:       "module.mymodule.resource.path[1]",
				ModuleAddress: "module.mymodule",
				Type:          "resource",
				Name:          "path",
				Index:         1,
				UpdateType:    NewResource,
			},
		},
		"handles modules with data": {
			line:        "    # module.mymodule.data.mydata.path will be read during apply",
			shouldError: false,
			expected: &ResourceChange{
				Address:       "module.mymodule.data.mydata.path",
				ModuleAddress: "module.mymodule",
				Type:          "mydata",
				Name:          "path",
				UpdateType:    ReadResource,
			},
		},
		"handles modules with data and data as the name": {
			line:        "    # module.mymodule.data.data.data will be read during apply",
			shouldError: false,
			expected: &ResourceChange{
				Address:       "module.mymodule.data.data.data",
				ModuleAddress: "module.mymodule",
				Type:          "data",
				Name:          "data",
				UpdateType:    ReadResource,
			},
		},
		"other line": {
			line:        "~ resource",
			shouldError: true,
			expected:    nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := NewResourceChangeFromComment(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if err != nil && !tc.shouldError {
				t.Fatalf("Unexpected error %v", err)
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestGetBeforeResource(t *testing.T) {
	cases := map[string]struct {
		rc       *ResourceChange
		expected map[string]interface{}
		opts     []GetBeforeAfterOptions
	}{
		"one attribute": {
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				MapAttributeChanges: []*MapAttributeChange{
					&MapAttributeChange{
						Name: "map",
						AttributeChanges: []*AttributeChange{
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
		"map and normal attribute": {
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
				},
				MapAttributeChanges: []*MapAttributeChange{
					&MapAttributeChange{
						Name: "map",
						AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			if got := tc.rc.GetBeforeResource(tc.opts...); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}

func TestGetAfterResource(t *testing.T) {
	cases := map[string]struct {
		rc       *ResourceChange
		expected map[string]interface{}
		opts     []GetBeforeAfterOptions
	}{
		"one attribute": {
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				MapAttributeChanges: []*MapAttributeChange{
					&MapAttributeChange{
						Name: "map",
						AttributeChanges: []*AttributeChange{
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
		"map and normal attribute": {
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
					&AttributeChange{
						Name:     "attribute",
						OldValue: "oldValue",
						NewValue: "newValue",
					},
				},
				MapAttributeChanges: []*MapAttributeChange{
					&MapAttributeChange{
						Name: "map",
						AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			rc: &ResourceChange{
				AttributeChanges: []*AttributeChange{
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
			if got := tc.rc.GetAfterResource(tc.opts...); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
