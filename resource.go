package tfplanparse

import (
	"fmt"
	"strconv"
	"strings"
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
	AttributeChanges []attributeChange
}

// IsResourceCommentLine returns true if the line is a valid resource comment line
// A valid line starts with a "#" and has a suffix describing the change
// Example: # module.type.item will be created
func IsResourceCommentLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "#") && !strings.HasSuffix(trimmed, RESOURCE_READ_VALUES_NOT_YET_KNOWN)
}

// IsResourceTerminator returns true if the line is a "}"
func IsResourceTerminator(line string) bool {
	return strings.TrimSpace(line) == "}"
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
	var address string

	// parse index first in case the index contains a "."
	addressIndex := strings.Split(rc.Address, "[")
	address = addressIndex[0]

	if len(addressIndex) == 2 {
		index := dequote(strings.TrimSuffix(addressIndex[1], "]"))

		if i, err := strconv.Atoi(index); err == nil {
			rc.Index = i
		} else {
			rc.Index = index
		}
	} else if len(addressIndex) > 2 {
		return fmt.Errorf("failed to parse resource info from address %s", rc.Address)
	}

	values := strings.Split(address, ".")

	// TODO: handle module.module_name.data.type.name better
	// TODO: eventually do something with "data"
	// For now, since we're not handling it, we can just remove it
	for k, v := range values {
		var previous string
		if k != 0 {
			previous = values[k-1]
		}

		// don't remove "data" if any of the conditions are true:
		// 1. Previous element was "module" or "data" (this means the module or data itself is named "data")
		// 2. There are less than 2 elements left to parse (this means the resource name or type is "data")
		if v == "data" && (previous != "module" && previous != "data") && (len(values)-k) > 2 {
			values = append(values[:k], values[k+1:]...)
		}
	}

	if len(values) == 2 {
		rc.Name = values[1]
		rc.Type = values[0]
	} else if len(values) > 2 {
		rc.Name = values[len(values)-1]
		rc.Type = values[len(values)-2]
		rc.ModuleAddress = fmt.Sprintf("%s.%s", values[0], values[1])
	} else {
		return fmt.Errorf("failed to parse resource info from address %s", rc.Address)
	}

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
		result[a.GetName()] = a.GetBefore(opts...)
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
		result[a.GetName()] = a.GetAfter(opts...)
	}

	return result
}

func parseResourceAddressFromComment(comment, updateText string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "# "), updateText))
}
