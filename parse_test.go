package tfplanparse

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		file     string
		expected []*ResourceChange
	}{
		// "basic plan": {
		// 	file: "test/basic.stdout",
		// 	expected: []*ResourceChange{
		// 		// TODO
		// 	},
		// },
		// "plan error": {
		// 	file: "test/error.stdout",
		// 	expected: nil,
		// },
		// "no changes": {
		// 	file: "test/nochanges.stdout",
		// 	expected: nil,
		// },
		// "array": {
		// 	file: "test/array.stdout",
		// 	expected: []*ResourceChange{
		// 		// TODO
		// 	},
		// },
		"nested map": {
			file: "test/nestedmap.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:          "module.mymodule.kubernetes_namespace.mynamespace",
					ModuleAddress:    "module.mymodule",
					Type:             "kubernetes_namespace",
					Name:             "mynamespace",
					UpdateType:       UpdateInPlaceResource,
					AttributeChanges: nil,
					MapAttributeChanges: []*MapAttributeChange{
						&MapAttributeChange{
							Name:             "metadata",
							AttributeChanges: nil,
							MapAttributeChanges: []*MapAttributeChange{
								&MapAttributeChange{
									Name: "labels",
									AttributeChanges: []*AttributeChange{
										&AttributeChange{
											Name:       "newLabel",
											OldValue:   nil,
											NewValue:   "newLabel",
											UpdateType: NewResource,
										},
									},
									MapAttributeChanges: nil,
									UpdateType:          UpdateInPlaceResource,
								},
							},
							UpdateType: UpdateInPlaceResource,
						},
					},
				},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// TODO handle expected error
			got, err := ParseFromFile(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(got, tc.expected); diff != "" {
				t.Errorf("(-got, +expected)\n%s", diff)
			}
		})
	}
}
