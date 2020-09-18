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
		"resources": {
			file: "test/resources.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:       "module.my-module.github_team_membership.member",
					ModuleAddress: "module.my-module",
					Type:          "github_team_membership",
					Name:          "member",
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "etag",
							OldValue:   `W/\"etag-0\"`,
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "id",
							OldValue:   "1234567:dev0",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "role",
							OldValue:   "member",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "team_id",
							OldValue:   "1234567",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "username",
							OldValue:   "dev0",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
					},
				},
				&ResourceChange{
					Address:       "module.my-module.github_team_membership.member[1]",
					ModuleAddress: "module.my-module",
					Type:          "github_team_membership",
					Name:          "member",
					Index:         1,
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "etag",
							OldValue:   `W/\"etag-1\"`,
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "id",
							OldValue:   "1234567:dev1",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "role",
							OldValue:   "member",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "team_id",
							OldValue:   "1234567",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "username",
							OldValue:   "dev1",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
					},
				},
				&ResourceChange{
					Address:       "module.my-module.github_team_membership.member[2]",
					ModuleAddress: "module.my-module",
					Type:          "github_team_membership",
					Name:          "member",
					Index:         2,
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "etag",
							OldValue:   `W/\"etag-2\"`,
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "id",
							OldValue:   "1234567:dev2",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "role",
							OldValue:   "member",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "team_id",
							OldValue:   "1234567",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "username",
							OldValue:   "dev2",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
					},
				},
				&ResourceChange{
					Address:       "module.my-module.github_team_membership.member[3]",
					ModuleAddress: "module.my-module",
					Type:          "github_team_membership",
					Name:          "member",
					Index:         3,
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "etag",
							OldValue:   `W/\"etag-3\"`,
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "id",
							OldValue:   "1234567:dev3",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "role",
							OldValue:   "member",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "team_id",
							OldValue:   "1234567",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "username",
							OldValue:   "dev3",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
					},
				},
			},
		},
		"array": {
			file: "test/array.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:       "module.my-project.google_project_services.gcp_enabled_services[0]",
					ModuleAddress: "module.my-project",
					Type:          "google_project_services",
					Name:          "gcp_enabled_services",
					Index:         0,
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "disable_on_destroy",
							OldValue:   false,
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "id",
							OldValue:   "my-project",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&AttributeChange{
							Name:       "project",
							OldValue:   "my-project",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&ArrayAttributeChange{
							Name: "services",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									OldValue:   "appengine.googleapis.com",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									OldValue:   "audit.googleapis.com",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name:       "timeouts",
							UpdateType: DestroyResource,
						},
					},
				},
			},
		},
		"nested map": {
			file: "test/nestedmap.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:       "module.mymodule.kubernetes_namespace.mynamespace",
					ModuleAddress: "module.mymodule",
					Type:          "kubernetes_namespace",
					Name:          "mynamespace",
					UpdateType:    UpdateInPlaceResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "id",
							OldValue:   "namespace-id",
							NewValue:   "namespace-id",
							UpdateType: NoOpResource,
						},
						&MapAttributeChange{
							Name: "metadata",
							AttributeChanges: []attributeChange{
								&MapAttributeChange{
									Name:       "annotations",
									UpdateType: NoOpResource,
								},
								&AttributeChange{
									Name:       "generation",
									OldValue:   0,
									NewValue:   0,
									UpdateType: NoOpResource,
								},
								&MapAttributeChange{
									Name: "labels",
									AttributeChanges: []attributeChange{
										&AttributeChange{
											Name:       "label",
											OldValue:   "value",
											NewValue:   "value",
											UpdateType: NoOpResource,
										},
										&AttributeChange{
											Name:       "other",
											OldValue:   "label",
											NewValue:   "label",
											UpdateType: NoOpResource,
										},
										&AttributeChange{
											Name:       "newLabel",
											OldValue:   nil,
											NewValue:   "newLabel",
											UpdateType: NewResource,
										},
									},
									UpdateType: UpdateInPlaceResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "my-namespace",
									NewValue:   "my-namespace",
									UpdateType: NoOpResource,
								},
								&AttributeChange{
									Name:       "resource_version",
									OldValue:   "123",
									NewValue:   "123",
									UpdateType: NoOpResource,
								},
								&AttributeChange{
									Name:       "self_link",
									OldValue:   "/api/v1/namespaces/my-namespace",
									NewValue:   "/api/v1/namespaces/my-namespace",
									UpdateType: NoOpResource,
								},
								&AttributeChange{
									Name:       "uid",
									OldValue:   "some-uid-123",
									NewValue:   "some-uid-123",
									UpdateType: NoOpResource,
								},
							},
							UpdateType: UpdateInPlaceResource,
						},
						&MapAttributeChange{
							Name:       "timeouts",
							UpdateType: NoOpResource,
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
