package tfplanparse

import (
	"fmt"
	"strings"
)

type MapAttributeChange struct {
	Name                    string
	AttributeChanges        []*AttributeChange
	MapAttributeChanges     []*MapAttributeChange
	ArrayAttributeChanges   []*ArrayAttributeChange
	HeredocAttributeChanges []*HeredocAttributeChange
	UpdateType              UpdateType
}

// IsMapAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", not be followed with "resource" or "data", and ends with "{".
func IsMapAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	// validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	validSuffix := strings.HasSuffix(line, "{") || IsOneLineEmptyMapAttribute(line)
	return validSuffix && !IsResourceChangeLine(line)
}

// IsMapAttributeTerminator returns true if the line is a "}" or "},"
func IsMapAttributeTerminator(line string) bool {
	return strings.TrimSuffix(strings.TrimSpace(line), ",") == "}"
}

// IsOneLineEmptyMapAttribute returns true if the line ends with a "{}"
func IsOneLineEmptyMapAttribute(line string) bool {
	return strings.HasSuffix(line, "{}")
}

// NewMapAttributeChangeFromLine initializes an AttributeChange from a line containing an attribute change
// It expects a line that passes the IsAttributeChangeLine check
func NewMapAttributeChangeFromLine(line string) (*MapAttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsMapAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a MapAttributeChange", line)
	}

	attributeName := getMultiLineAttributeName(line)
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
		return &MapAttributeChange{
			Name:       attributeName,
			UpdateType: NoOpResource,
		}, nil
	}
}

func (m *MapAttributeChange) GetBeforeAttribute(opts ...GetBeforeAfterOptions) map[string]interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.Name] = a.OldValue
	}

	for _, ma := range m.MapAttributeChanges {
		result[ma.Name] = ma.GetBeforeAttribute(opts...)
	}

	for _, aa := range m.ArrayAttributeChanges {
		result[aa.Name] = aa.GetBeforeAttribute(opts...)
	}

	for _, ha := range m.HeredocAttributeChanges {
		result[ha.Name] = ha.GetBeforeAttribute(opts...)
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

	for _, aa := range m.ArrayAttributeChanges {
		result[aa.Name] = aa.GetAfterAttribute(opts...)
	}

	for _, ha := range m.HeredocAttributeChanges {
		result[ha.Name] = ha.GetAfterAttribute(opts...)
	}

	return result
}
