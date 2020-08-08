package tfplanparse

import (
	"fmt"
	"strings"
)

type ResourcePlan struct {
	Path                string
	Type                string
	Name                string
	UpdateType          UpdateType
	Tainted             bool
	AttributeChanges    []*AttributeChange
	MapAttributeChanges []*MapAttributeChange
}

const (
	RESOURCE_CREATED                   = " will be created"
	RESOURCE_READ                      = " will be read during apply"
	RESOURCE_READ_VALUES_NOT_YET_KNOWN = " (config refers to values not yet known)"
	RESOURCE_UPDATED_IN_PLACE          = " will be updated in-place"
	RESOURCE_TAINTED                   = " is tainted, so must be replaced"
	RESOURCE_REPLACED                  = " must be replaced"
	RESOURCE_DESTROYED                 = " will be destroyed"
)

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

// NewResourcePlanFromComment creates a ResourcePlan from a valid resource comment line
func NewResourcePlanFromComment(comment string) (*ResourcePlan, error) {
	comment = strings.TrimSpace(comment)
	if !IsResourceCommentLine(comment) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a resource", comment)
	}

	if strings.HasSuffix(comment, RESOURCE_CREATED) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "# "), RESOURCE_CREATED))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: NewResource,
		}, nil
	} else if strings.HasSuffix(comment, RESOURCE_READ) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "#"), RESOURCE_READ))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: ReadResource,
		}, nil
	} else if strings.HasSuffix(comment, RESOURCE_UPDATED_IN_PLACE) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "#"), RESOURCE_UPDATED_IN_PLACE))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: UpdateInPlaceResource,
		}, nil
	} else if strings.HasSuffix(comment, RESOURCE_TAINTED) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "#"), RESOURCE_TAINTED))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: ForceReplaceResource,
			Tainted:    true,
		}, nil
	} else if strings.HasSuffix(comment, RESOURCE_REPLACED) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "#"), RESOURCE_REPLACED))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: ForceReplaceResource,
		}, nil
	} else if strings.HasSuffix(comment, RESOURCE_DESTROYED) {
		resourcePath := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(comment, "#"), RESOURCE_DESTROYED))
		resourceType, resourceName := parseResourceTypeAndName(resourcePath)
		return &ResourcePlan{
			Path:       resourcePath,
			Type:       resourceType,
			Name:       resourceName,
			UpdateType: DestroyResource,
		}, nil
	}

	return nil, fmt.Errorf("unknown comment line %s", comment)
}

func parseResourceTypeAndName(line string) (string, string) {
	values := strings.Split(line, ".")
	if len(values) == 2 {
		return values[0], values[1]
	} else if len(values) == 4 {
		return values[2], values[3]
	}

	return "UNKNOWN_TYPE", "UNKNOWN_NAME"
}
