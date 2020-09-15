package tfplanparse

import (
	"fmt"
	"strings"
)

// TODO: array attributes can be any attribute change, an array of arrays, or array of maps
// Type of the attribute can vary, but they're all the same
// Should use an interface
type ArrayAttributeChange struct {
	Name                  string
	AttributeChanges      []*AttributeChange
	MapAttributeChanges   []*MapAttributeChange
	ArrayAttributeChanges []*ArrayAttributeChange
	UpdateType            UpdateType
}

// IsArrayAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", not be followed with "resource" or "data", and ends with "[".
func IsArrayAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	validSuffix := strings.HasSuffix(line, "[") || IsOneLineEmptyArrayAttribute(line)
	return validPrefix && validSuffix && !IsResourceChangeLine(line)
}

// IsArrayAttributeTerminator returns true if the line is "]" or "] -> null"
func IsArrayAttributeTerminator(line string) bool {
	return strings.TrimSuffix(strings.TrimSpace(line), " -> null") == "]"
}

// IsOneLineEmptyArrayAttribute returns true if the line ends with a "[]"
func IsOneLineEmptyArrayAttribute(line string) bool {
	return strings.HasSuffix(line, "[]")
}

// NewArrayAttributeChangeFromLine initializes an ArrayAttributeChange from a line containing an array attribute change
// It expects a line that passes the IsArrayAttributeChangeLine check
func NewArrayAttributeChangeFromLine(line string) (*ArrayAttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsArrayAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a ArrayAttributeChange", line)
	}

	attributeName := getMultiLineAttributeName(line)
	if strings.HasPrefix(line, "+") {
		// add
		return &ArrayAttributeChange{
			Name:       attributeName,
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		return &ArrayAttributeChange{
			Name:       attributeName,
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		return &ArrayAttributeChange{
			Name:       attributeName,
			UpdateType: UpdateInPlaceResource,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern")
	}
}

func (m *ArrayAttributeChange) GetBeforeAttribute(opts ...GetBeforeAfterOptions) []interface{} {
	// TODO: ensure the result types are all the same
	// Currently it is assumed that all changes added are the same type...
	// This is handled correctly in parse, but we should handle it here
	result := []interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		if a.UpdateType == NewResource {
			continue attrs
		}
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result = append(result, a.OldValue)
	}

	for _, aa := range m.ArrayAttributeChanges {
		result = append(result, aa.GetBeforeAttribute(opts...))
	}

	for _, ma := range m.MapAttributeChanges {
		result = append(result, ma.GetBeforeAttribute(opts...))
	}

	return result
}

func (m *ArrayAttributeChange) GetAfterAttribute(opts ...GetBeforeAfterOptions) []interface{} {
	// TODO: same as above
	result := []interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		if a.UpdateType == DestroyResource {
			continue attrs
		}
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result = append(result, a.NewValue)
	}

	for _, aa := range m.ArrayAttributeChanges {
		result = append(result, aa.GetAfterAttribute(opts...))
	}

	for _, ma := range m.MapAttributeChanges {
		result = append(result, ma.GetAfterAttribute(opts...))
	}

	return result
}
