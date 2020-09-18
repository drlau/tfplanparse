package tfplanparse

import (
	"fmt"
	"strings"
)

// TODO: array attributes can be any attribute change, an array of arrays, or array of maps
// Type of the attribute can vary, but they're all the same
// Should use an interface
type ArrayAttributeChange struct {
	Name             string
	AttributeChanges []attributeChange
	UpdateType       UpdateType
}

var _ attributeChange = &ArrayAttributeChange{}

// IsArrayAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", not be followed with "resource" or "data", and ends with "[".
func IsArrayAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	// validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	validSuffix := strings.HasSuffix(line, "[") || IsOneLineEmptyArrayAttribute(line)
	return validSuffix && !IsResourceChangeLine(line)
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
		return &ArrayAttributeChange{
			Name:       attributeName,
			UpdateType: NoOpResource,
		}, nil
	}
}

// GetName returns the name of the attribute
func (a *ArrayAttributeChange) GetName() string {
	return a.Name
}

// GetUpdateType returns the UpdateType of the attribute
func (a *ArrayAttributeChange) GetUpdateType() UpdateType {
	return a.UpdateType
}

// IsSensitive returns true if the attribute contains a sensitive value
func (a *ArrayAttributeChange) IsSensitive() bool {
	// return m.OldValue == SENSITIVE_VALUE || m.NewValue == SENSITIVE_VALUE
	return false
}

// IsComputed returns true if the attribute contains a computed value
func (a *ArrayAttributeChange) IsComputed() bool {
	// return m.OldValue == COMPUTED_VALUE || m.NewValue == COMPUTED_VALUE
	return false
}

// IsNoOp returns true if the attribute has not changed
func (a *ArrayAttributeChange) IsNoOp() bool {
	return a.UpdateType == NoOpResource
}

func (a *ArrayAttributeChange) GetBefore(opts ...GetBeforeAfterOptions) interface{} {
	// TODO: ensure the result types are all the same
	// Currently it is assumed that all changes added are the same type...
	// This is handled correctly in parse, but we should handle it here
	result := []interface{}{}

attrs:
	for _, ac := range a.AttributeChanges {
		if ac.GetUpdateType() == NewResource {
			continue attrs
		}
		for _, opt := range opts {
			if opt(ac) {
				continue attrs
			}
		}
		result = append(result, ac.GetBefore(opts...))
	}

	return result
}

func (a *ArrayAttributeChange) GetAfter(opts ...GetBeforeAfterOptions) interface{} {
	// TODO: same as above
	result := []interface{}{}

attrs:
	for _, ac := range a.AttributeChanges {
		if ac.GetUpdateType() == DestroyResource {
			continue attrs
		}
		for _, opt := range opts {
			if opt(ac) {
				continue attrs
			}
		}
		result = append(result, ac.GetAfter(opts...))
	}

	return result
}
