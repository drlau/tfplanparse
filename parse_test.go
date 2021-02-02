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
		"another map": {
			file: "test/anothermap.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:       "module.mymodule.kubernetes_role_binding.user_is_edit",
					ModuleAddress: "module.mymodule",
					Type:          "kubernetes_role_binding",
					Name:          "user_is_edit",
					UpdateType:    NewResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "id",
							OldValue:   nil,
							NewValue:   "(known after apply)",
							UpdateType: NewResource,
						},
						&MapAttributeChange{
							Name: "metadata",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "generation",
									OldValue:   nil,
									NewValue:   "(known after apply)",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   nil,
									NewValue:   "user-is-edit",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "namespace",
									OldValue:   nil,
									NewValue:   "my-namespace",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "resource_version",
									OldValue:   nil,
									NewValue:   "(known after apply)",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "self_link",
									OldValue:   nil,
									NewValue:   "(known after apply)",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "uid",
									OldValue:   nil,
									NewValue:   "(known after apply)",
									UpdateType: NewResource,
								},
							},
							UpdateType: NewResource,
						},
						&MapAttributeChange{
							Name: "role_ref",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   nil,
									NewValue:   "rbac.authorization.k8s.io",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   nil,
									NewValue:   "ClusterRole",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   nil,
									NewValue:   "edit",
									UpdateType: NewResource,
								},
							},
							UpdateType: NewResource,
						},
						&MapAttributeChange{
							Name: "subject",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   nil,
									NewValue:   "rbac.authorization.k8s.io",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   nil,
									NewValue:   "User",
									UpdateType: NewResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   nil,
									NewValue:   "user@email.com",
									UpdateType: NewResource,
								},
							},
							UpdateType: NewResource,
						},
					},
				},
				&ResourceChange{
					Address:       "module.mymodule.kubernetes_role_binding.user_is_view",
					ModuleAddress: "module.mymodule",
					Type:          "kubernetes_role_binding",
					Name:          "user_is_view",
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "id",
							OldValue:   "my-namespace/user_is_view",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "metadata",
							AttributeChanges: []attributeChange{
								&MapAttributeChange{
									Name: "annotations",
									AttributeChanges: []attributeChange{
										&AttributeChange{
											Name:       "my-annotation",
											OldValue:   "annot",
											NewValue:   nil,
											UpdateType: DestroyResource,
										},
									},
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "generation",
									OldValue:   0,
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "labels",
									OldValue:   nil,
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "user-is-view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "namespace",
									OldValue:   "my-namespace",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "resource_version",
									OldValue:   "123",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "self_link",
									OldValue:   "/apis/rbac.authorization.k8s.io/v1/namespaces/my-namespace/rolebindings/user-is-view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "uid",
									OldValue:   "some-uid",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "role_ref",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   "rbac.authorization.k8s.io",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   "ClusterRole",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "subject",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   "rbac.authorization.k8s.io",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   "User",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "user@email.com",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
					},
				},
			},
		},
		"jsonencode": {
			file: "test/jsonencode.stdout",
			expected: []*ResourceChange{
				&ResourceChange{
					Address:       "module.mymodule.kubernetes_role_binding.user_is_view",
					ModuleAddress: "module.mymodule",
					Type:          "kubernetes_role_binding",
					Name:          "user_is_view",
					UpdateType:    DestroyResource,
					AttributeChanges: []attributeChange{
						&AttributeChange{
							Name:       "id",
							OldValue:   "my-namespace/user_is_view",
							NewValue:   nil,
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "metadata",
							AttributeChanges: []attributeChange{
								&MapAttributeChange{
									Name: "annotations",
									AttributeChanges: []attributeChange{
										&JSONEncodeAttributeChange{
											Name: "encoded",
											AttributeChanges: []attributeChange{
												&MapAttributeChange{
													AttributeChanges: []attributeChange{
														&AttributeChange{
															Name:       "apiVersion",
															OldValue:   "rbac.authorization.k8s.io/v1",
															NewValue:   nil,
															UpdateType: DestroyResource,
														},
														&AttributeChange{
															Name:       "kind",
															OldValue:   "RoleBinding",
															NewValue:   nil,
															UpdateType: DestroyResource,
														},
														&MapAttributeChange{
															Name: "metadata",
															AttributeChanges: []attributeChange{
																&MapAttributeChange{
																	Name: "annotations",
																	AttributeChanges: []attributeChange{
																		&AttributeChange{
																			Name:       "my-annotation",
																			OldValue:   "annot",
																			NewValue:   nil,
																			UpdateType: DestroyResource,
																		},
																	},
																	UpdateType: DestroyResource,
																},
																&AttributeChange{
																	Name:       "creationTimestamp",
																	OldValue:   "null",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
																&MapAttributeChange{
																	Name: "labels",
																	AttributeChanges: []attributeChange{
																		&AttributeChange{
																			Name:       "my-label",
																			OldValue:   "label",
																			NewValue:   nil,
																			UpdateType: DestroyResource,
																		},
																	},
																	UpdateType: DestroyResource,
																},
																&AttributeChange{
																	Name:       "name",
																	OldValue:   "user-is-view",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
																&AttributeChange{
																	Name:       "namespace",
																	OldValue:   "my-namespace",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
															},
															UpdateType: DestroyResource,
														},
														&MapAttributeChange{
															Name: "roleRef",
															AttributeChanges: []attributeChange{
																&AttributeChange{
																	Name:       "apiGroup",
																	OldValue:   "rbac.authorization.k8s.io",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
																&AttributeChange{
																	Name:       "kind",
																	OldValue:   "ClusterRole",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
																&AttributeChange{
																	Name:       "name",
																	OldValue:   "view",
																	NewValue:   nil,
																	UpdateType: DestroyResource,
																},
															},
															UpdateType: DestroyResource,
														},
														&ArrayAttributeChange{
															Name: "subjects",
															AttributeChanges: []attributeChange{
																&MapAttributeChange{
																	Name: "",
																	AttributeChanges: []attributeChange{
																		&AttributeChange{
																			Name:       "apiGroup",
																			OldValue:   "rbac.authorization.k8s.io",
																			NewValue:   nil,
																			UpdateType: DestroyResource,
																		},
																		&AttributeChange{
																			Name:       "kind",
																			OldValue:   "User",
																			NewValue:   nil,
																			UpdateType: DestroyResource,
																		},
																		&AttributeChange{
																			Name:       "name",
																			OldValue:   "user@email.com",
																			NewValue:   nil,
																			UpdateType: DestroyResource,
																		},
																	},
																	UpdateType: DestroyResource,
																},
															},
															UpdateType: DestroyResource,
														},
													},
													// TODO: this should be DestroyResource
													UpdateType: NoOpResource,
												},
											},
											UpdateType: DestroyResource,
										},
									},
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "generation",
									OldValue:   0,
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "labels",
									OldValue:   nil,
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "user-is-view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "namespace",
									OldValue:   "my-namespace",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "resource_version",
									OldValue:   "123",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "self_link",
									OldValue:   "/apis/rbac.authorization.k8s.io/v1/namespaces/my-namespace/rolebindings/user-is-view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "uid",
									OldValue:   "some-uid",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "role_ref",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   "rbac.authorization.k8s.io",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   "ClusterRole",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "view",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
						},
						&MapAttributeChange{
							Name: "subject",
							AttributeChanges: []attributeChange{
								&AttributeChange{
									Name:       "api_group",
									OldValue:   "rbac.authorization.k8s.io",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "kind",
									OldValue:   "User",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
								&AttributeChange{
									Name:       "name",
									OldValue:   "user@email.com",
									NewValue:   nil,
									UpdateType: DestroyResource,
								},
							},
							UpdateType: DestroyResource,
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
