package tfplanparse

import (
	"fmt"
	"strings"
)

type JSONEncodeAttributeChange struct {
	Name             string
	AttributeChanges []attributeChange
	UpdateType       UpdateType
}

var _ attributeChange = &JSONEncodeAttributeChange{}

// IsJSONEncodeAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", delimited with a space, and the value to start with "jsonencode(".
func IsJSONEncodeAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	attribute := strings.SplitN(line, ATTRIBUTE_DEFINITON_DELIMITER, 2)
	if len(attribute) != 2 {
		return false
	}

	validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	isJSONEncode := strings.HasPrefix(attribute[1], "jsonencode(")
	return validPrefix && isJSONEncode && !IsResourceChangeLine(line)
}

// IsJSONEncodeAttributeTerminator returns true if the line is ")"
// TODO: verify this
func IsJSONEncodeAttributeTerminator(line string) bool {
	return strings.TrimSuffix(strings.TrimSpace(line), " -> null") == ")"
}

// NewJSONEncodeAttributeChangeFromLine initializes a JSONEncodeAttributeChange from a line containing a JSONEncode change
// It expects a line that passes the IsJSONEncodeAttributeChangeLine check
func NewJSONEncodeAttributeChangeFromLine(line string) (*JSONEncodeAttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsJSONEncodeAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a JSONEncodeAttributeChange", line)
	}
	attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)

	if strings.HasPrefix(line, "+") {
		// add
		return &JSONEncodeAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		return &JSONEncodeAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		updateType := UpdateInPlaceResource
		if strings.HasSuffix(attribute[1], " # forces replacement") {
			updateType = ForceReplaceResource
		}

		return &JSONEncodeAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			UpdateType: updateType,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern")
	}
}

// GetName returns the name of the attribute
func (j *JSONEncodeAttributeChange) GetName() string {
	return j.Name
}

// GetUpdateType returns the UpdateType of the attribute
func (j *JSONEncodeAttributeChange) GetUpdateType() UpdateType {
	return j.UpdateType
}

// IsSensitive returns true if the attribute contains a sensitive value
func (j *JSONEncodeAttributeChange) IsSensitive() bool {
	for _, ac := range j.AttributeChanges {
		if ac.IsSensitive() {
			return true
		}
	}
	return false
}

// IsComputed returns true if the attribute contains a computed value
func (j *JSONEncodeAttributeChange) IsComputed() bool {
	for _, ac := range j.AttributeChanges {
		if ac.IsComputed() {
			return true
		}
	}
	return false
}

// IsNoOp returns true if the attribute has not changed
func (j *JSONEncodeAttributeChange) IsNoOp() bool {
	return j.UpdateType == NoOpResource
}

func (j *JSONEncodeAttributeChange) GetBefore(opts ...GetBeforeAfterOptions) interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range j.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.GetName()] = a.GetBefore(opts...)
	}

	return result
}

func (j *JSONEncodeAttributeChange) GetAfter(opts ...GetBeforeAfterOptions) interface{} {
	result := map[string]interface{}{}

attrs:
	for _, a := range j.AttributeChanges {
		for _, opt := range opts {
			if opt(a) {
				continue attrs
			}
		}
		result[a.GetName()] = a.GetAfter(opts...)
	}

	return result
}
