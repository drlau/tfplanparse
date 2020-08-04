package plan

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

func TestNewResourcePlanFromComment(t *testing.T) {
	cases := map[string]struct {
		line        string
		shouldError bool
		expected    *ResourcePlan
	}{
		"empty line": {
			line:        "",
			shouldError: true,
			expected:    nil,
		},
		"resource created": {
			line:        "# resource.path will be created",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: NewResource,
			},
		},
		"resource read during apply": {
			line:        "# resource.path will be read during apply",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
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
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: UpdateInPlaceResource,
			},
		},
		"resource tainted": {
			line:        "# resource.path is tainted, so must be replaced",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ForceReplaceResource,
				Tainted:    true,
			},
		},
		"resource replaced": {
			line:        "# resource.path must be replaced",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ForceReplaceResource,
			},
		},
		"resource destroyed": {
			line:        "# resource.path will be destroyed",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: DestroyResource,
			},
		},
		"handles extra spaces": {
			line:        "    # resource.path will be created",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: NewResource,
			},
		},
		"handles data paths spaces": {
			line:        "    # data.mydata.resource.path will be read during apply",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "data.mydata.resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: ReadResource,
			},
		},
		"handles modules with extra spaces": {
			line:        "    # module.mymodule.resource.path will be created",
			shouldError: false,
			expected: &ResourcePlan{
				Path:       "module.mymodule.resource.path",
				Type:       "resource",
				Name:       "path",
				UpdateType: NewResource,
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
			got, err := NewResourcePlanFromComment(tc.line)
			if err == nil && tc.shouldError {
				t.Fatalf("Expected an error but didn't get one")
			}

			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("Expected: %v but got %v", tc.expected, got)
			}
		})
	}
}
