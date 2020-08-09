package tfplanparse

import (
	"fmt"
	"strings"
)

type MapAttributeChange struct {
	Name                string
	AttributeChanges    []*AttributeChange
	MapAttributeChanges []*MapAttributeChange
	UpdateType          UpdateType
}

// IsMapAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", not be followed with "resource" or "data", and ends with "{".
func IsMapAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	validSuffix := strings.HasSuffix(line, "{")
	return validPrefix && validSuffix && !IsResourceChangeLine(line)
}

// IsMapAttributeTerminator returns true if the line is a "}"
func IsMapAttributeTerminator(line string) bool {
	return strings.TrimSpace(line) == "}"
}

// NewMapAttributeChangeFromLine initializes an AttributeChange from a line containing an attribute change
// It expects a line that passes the IsAttributeChangeLine check
func NewMapAttributeChangeFromLine(line string) (*MapAttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsMapAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a MapAttributeChange", line)
	}

	attributeName := getMapAttributeName(line)
	if strings.HasPrefix(line, "+") {
		// add
		return &MapAttributeChange{
			Name:       attributeName,
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		return &MapAttributeChange{
			Name:       attributeName,
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		return &MapAttributeChange{
			Name:       attributeName,
			UpdateType: UpdateInPlaceResource,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern")
	}
}

func (m *MapAttributeChange) GetBeforeAttribute() map[string]interface{} {
	result := map[string]interface{}{}

	for _, a := range m.AttributeChanges {
		result[a.Name] = a.OldValue
	}

	for _, ma := range m.MapAttributeChanges {
		result[ma.Name] = ma.GetBeforeAttribute()
	}

	return result
}

func (m *MapAttributeChange) GetAfterAttribute(opts ...GetBeforeAfterOptions) map[string]interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.Name] = a.NewValue
	}

	for _, ma := range m.MapAttributeChanges {
		result[ma.Name] = ma.GetAfterAttribute(opts...)
	}

	return result
}

func getMapAttributeName(line string) string {
	line = removeChangeTypeCharacters(line)
	// Map attributes may or may not have a name
	// If they do have a name, they are delimited with a '=' or a ' '
	attribute := strings.SplitN(line, ATTRIBUTE_DEFINITON_DELIMITER, 2)
	if len(attribute) == 2 {
		return strings.TrimSpace(attribute[0])
	}

	attribute = strings.SplitN(line, " ", 2)
	if len(attribute) == 2 {
		return strings.TrimSpace(attribute[0])
	}

	return ""
}
