package tfplanparse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-terraform-address"
)

const (
	RESOURCE_CREATED                   = " will be created"
	RESOURCE_READ                      = " will be read during apply"
	RESOURCE_READ_VALUES_NOT_YET_KNOWN = " (config refers to values not yet known)"
	RESOURCE_UPDATED_IN_PLACE          = " will be updated in-place"
	RESOURCE_TAINTED                   = " is tainted, so must be replaced"
	RESOURCE_REPLACED                  = " must be replaced"
	RESOURCE_DESTROYED                 = " will be destroyed"
)

type ResourceChange struct {
	// Address contains the absolute resource address
	Address string

	// ModuleAddress contains the module portion of the absolute address, if any
	ModuleAddress string

	// The type of the resource
	// Example: gcp_instance.foo -> "gcp_instance"
	Type string

	// The name of the resource
	// Example: gcp_instance.foo -> "foo"
	Name string

	// The index key for resources created with "count" or "for_each"
	// "count" resources will be an int index, and "for_each" will be a string
	Index interface{}

	// UpdateType contains the type of update
	// Refer to updatetype.go for possible values
	UpdateType UpdateType

	// Tainted indicates whether the resource is tainted or not
	Tainted bool

	// AttributeChanges contains all the planned attribute changes
	AttributeChanges []*AttributeChange

	// MapAttributeChanges contains all the planned attribute changes that are map type attributes
	MapAttributeChanges []*MapAttributeChange
}

// IsResourceCommentLine returns true if the line is a valid resource comment line
// A valid line starts with a "#" and has a suffix describing the change
// Example: # module.type.item will be created
func IsResourceCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#") && !strings.HasSuffix(trimmed, RESOURCE_READ_VALUES_NOT_YET_KNOWN)
}

// IsResourceChangeLine returns true if the line is a valid resource change line
// A valid line starts with the change type, then "resource" or "data", and then the type and name, followed by a {
// Example: + resource "type" "name" {
func IsResourceChangeLine(line string) bool {
	line = strings.TrimLeft(line, "+/-~<= ")
	return (strings.HasPrefix(line, "resource") || strings.HasPrefix(line, "data")) && strings.HasSuffix(line, " {")
}

// NewResourceChangeFromComment creates a ResourceChange from a valid resource comment line
func NewResourceChangeFromComment(comment string) (*ResourceChange, error) {
	var rc *ResourceChange
	var resourceAddress string
	comment = strings.TrimSpace(comment)
	if !IsResourceCommentLine(comment) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a resource", comment)
	}

	if strings.HasSuffix(comment, RESOURCE_CREATED) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_CREATED)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: NewResource,
		}
	} else if strings.HasSuffix(comment, RESOURCE_READ) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_READ)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: ReadResource,
		}
	} else if strings.HasSuffix(comment, RESOURCE_UPDATED_IN_PLACE) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_UPDATED_IN_PLACE)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: UpdateInPlaceResource,
		}
	} else if strings.HasSuffix(comment, RESOURCE_TAINTED) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_TAINTED)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: ForceReplaceResource,
			Tainted:    true,
		}
	} else if strings.HasSuffix(comment, RESOURCE_REPLACED) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_REPLACED)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: ForceReplaceResource,
		}
	} else if strings.HasSuffix(comment, RESOURCE_DESTROYED) {
		resourceAddress = parseResourceAddressFromComment(comment, RESOURCE_DESTROYED)

		rc = &ResourceChange{
			Address:    resourceAddress,
			UpdateType: DestroyResource,
		}
	}

	if rc == nil {
		return nil, fmt.Errorf("unknown comment line %s", comment)
	}

	if err := rc.finalizeResourceInfo(); err != nil {
		return nil, err
	}

	return rc, nil
}

func (rc *ResourceChange) finalizeResourceInfo() error {
	parsedInterface, err := address.Parse("", []byte(rc.Address))
	if err != nil {
		return err
	}

	parsed := parsedInterface.(*address.Address)

	rc.Type = parsed.ResourceSpec.Type
	rc.Name = parsed.ResourceSpec.Name
	rc.ModuleAddress = parsed.ModulePath.String()
	rc.Index = parsed.ResourceSpec.Index.Value

	return nil
}

func (rc *ResourceChange) GetBeforeResource(opts ...GetBeforeAfterOptions) map[string]interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range rc.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.Name] = a.OldValue
	}

	for _, m := range rc.MapAttributeChanges {
		result[m.Name] = m.GetBeforeAttribute(opts...)
	}

	return result
}

func (rc *ResourceChange) GetAfterResource(opts ...GetBeforeAfterOptions) map[string]interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range rc.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.Name] = a.NewValue
	}

	for _, m := range rc.MapAttributeChanges {
		result[m.Name] = m.GetAfterAttribute(opts...)
	}

	return result
}

func parseResourceAddressFromComment(comment, updateText string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "# "), updateText))
}
