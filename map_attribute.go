package tfplanparse

import (
	"fmt"
	"strings"
)

type MapAttributeChange struct {
	Name             string
	AttributeChanges []attributeChange
	UpdateType       UpdateType
}

var _ attributeChange = &MapAttributeChange{}

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

// GetName returns the name of the attribute
func (m *MapAttributeChange) GetName() string {
	return m.Name
}

// GetUpdateType returns the UpdateType of the attribute
func (m *MapAttributeChange) GetUpdateType() UpdateType {
	return m.UpdateType
}

// IsSensitive returns true if the attribute contains a sensitive value
func (m *MapAttributeChange) IsSensitive() bool {
	for _, ac := range m.AttributeChanges {
		if ac.IsSensitive() {
			return true
		}
	}
	return false
}

// IsComputed returns true if the attribute contains a computed value
func (m *MapAttributeChange) IsComputed() bool {
	for _, ac := range m.AttributeChanges {
		if ac.IsComputed() {
			return true
		}
	}
	return false
}

// IsNoOp returns true if the attribute has not changed
func (m *MapAttributeChange) IsNoOp() bool {
	return m.UpdateType == NoOpResource
}

func (m *MapAttributeChange) GetBefore(opts ...GetBeforeAfterOptions) interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.GetName()] = a.GetBefore(opts...)
	}

	return result
}

func (m *MapAttributeChange) GetAfter(opts ...GetBeforeAfterOptions) interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range m.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.GetName()] = a.GetAfter(opts...)
	}

	return result
}
