package tfplanparse

import (
	"fmt"
	"strings"
)

const (
	ATTRIBUTE_CHANGE_DELIMITER    = " -> "
	ATTRIBUTE_DEFINITON_DELIMITER = " = "
)

type AttributeChange struct {
	Name       string
	OldValue   string
	NewValue   string
	UpdateType UpdateType
}

// IsAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", and not be followed with "resource"
func IsAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	multilineAttribute := strings.HasSuffix(line, "(") || strings.HasSuffix(line, "{")

	return validPrefix && !multilineAttribute && !IsResourceChangeLine(line)
}

// NewAttributeChangeFromLine initializes an AttributeChange from a line containing an attribute change
// It expects a line that passes the IsAttributeChangeLine check
func NewAttributeChangeFromLine(line string) (*AttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize an attributeChange", line)
	}

	if strings.HasPrefix(line, "+") {
		// add
		attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)

		return &AttributeChange{
			Name:       strings.TrimSpace(attribute[0]),
			OldValue:   "",
			NewValue:   strings.TrimSpace(dequote(attribute[1])),
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)
		if len(attribute) == 1 {
			return &AttributeChange{
				Name:       strings.TrimSpace(attribute[0]),
				OldValue:   "",
				NewValue:   "",
				UpdateType: DestroyResource,
			}, nil
		}

		values := strings.Split(attribute[1], ATTRIBUTE_CHANGE_DELIMITER)
		if len(values) != 2 {
			return &AttributeChange{
				Name:       strings.TrimSpace(attribute[0]),
				OldValue:   strings.TrimSpace(attribute[1]),
				NewValue:   "",
				UpdateType: DestroyResource,
			}, nil
		}

		return &AttributeChange{
			Name:       strings.TrimSpace(attribute[0]),
			OldValue:   strings.TrimSpace(dequote(values[0])),
			NewValue:   "",
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		updateType := UpdateInPlaceResource

		if strings.HasSuffix(line, " # forces replacement") {
			updateType = ForceReplaceResource
			line = strings.TrimSuffix(line, " # forces replacement")
		}

		attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)

		values := strings.Split(attribute[1], ATTRIBUTE_CHANGE_DELIMITER)
		if len(values) != 2 {
			return nil, fmt.Errorf("failed to read attribute change")
		}

		return &AttributeChange{
			Name:       strings.TrimSpace(attribute[0]),
			OldValue:   strings.TrimSpace(dequote(values[0])),
			NewValue:   strings.TrimSpace(dequote(values[1])),
			UpdateType: updateType,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern")
	}
}

// IsSensitive returns true if the attribute contains a sensitive value
func (a *AttributeChange) IsSensitive() bool {
	return a.OldValue == "(sensitive value)" || a.NewValue == "(sensitive value)"
}

// IsComputed returns true if the attribute contains a computed value
func (a *AttributeChange) IsComputed() bool {
	return a.OldValue == "(known after apply)" || a.NewValue == "(known after apply)"
}

func removeChangeTypeCharacters(line string) string {
	return strings.TrimLeft(line, "+/-~<= ")
}

func dequote(line string) string {
	return strings.Trim(line, "\"")
}
