package tfplanparse

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
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

type ResourceIndex struct {
	Index interface{}
	Address string
}

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
	// This will only contain the final index, if you need the full index path, look at FullIndex
	Index interface{}

	// The full index path for resources created with "count" or "for_each"
	FullIndex []ResourceIndex

	// UpdateType contains the type of update
	// Refer to updatetype.go for possible values
	UpdateType UpdateType

	// Tainted indicates whether the resource is tainted or not
	Tainted bool

	// AttributeChanges contains all the planned attribute changes
	AttributeChanges []attributeChange
}

// IsComment returns true if the line is a comment.
//
// A valid line starts with a "#".
func IsComment(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

// IsResourceCommentLine returns true if the line is a resource "header" comment line
// A valid line starts with a "#" and ends with a known resource change suffix
// Example: # module.type.item will be created
func IsResourceCommentLine(line string) bool {
	line = strings.TrimSpace(line)

	if !IsComment(line) {
		return false
	}

	for _, suffix := range resourceCommentSuffixes {
		if strings.HasSuffix(line, suffix) {
			return true
		}
	}

	return false
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

// extractIndexes Extracts the final index from a resource path, along with any indexes on parent resources
func extractIndexes(address string) []ResourceIndex {
  getAllIndexs := regexp.MustCompile(`([^\[]*)(\[([^\]]*)\])?`)
	matches := getAllIndexs.FindAllStringSubmatch(address, -1)
	allIndexes := []ResourceIndex{}
	//fmt.Println("matches", address, matches)

	hasIndexes := false
	for _, match := range matches {
		var realIndex interface{}
		if len(match) > 2 && match[3] != "" {
			if i, err := strconv.Atoi(match[3]); err == nil {
				realIndex = i
			} else {
				realIndex = strings.Trim(strings.Trim(match[3], "\""), "'")
			}
			hasIndexes = true
		}
		allIndexes = append(allIndexes, ResourceIndex{
			Index: realIndex,
			Address: strings.TrimPrefix(match[1], "."),
		})
	}
	if !hasIndexes {
		return []ResourceIndex{}
	}
	return allIndexes
}

func (rc *ResourceChange) finalizeResourceInfo() error {
	// parse index first in case the index contains a "."
	allIndexes := extractIndexes(rc.Address)
	values := []string{}
	if len(allIndexes) > 0 {
		rc.Index = allIndexes[len(allIndexes)-1].Index
		rc.FullIndex = allIndexes
		for _, index := range allIndexes {
			for _, part := range strings.Split(index.Address, ".") {
				values = append(values, strings.TrimSuffix(strings.TrimPrefix(part, "."), "."))
			}
		}
	} else {
		values = strings.Split(rc.Address, ".")
	}

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

var resourceCommentSuffixes = []string{
	RESOURCE_CREATED, RESOURCE_DESTROYED, RESOURCE_READ,
	RESOURCE_REPLACED, RESOURCE_TAINTED, RESOURCE_UPDATED_IN_PLACE,
}
