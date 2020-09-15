package tfplanparse

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	ATTRIBUTE_CHANGE_DELIMITER    = " -> "
	ATTRIBUTE_DEFINITON_DELIMITER = " = "
	SENSITIVE_VALUE               = "(sensitive value)"
	COMPUTED_VALUE                = "(known after apply)"
)

type AttributeChange struct {
	Name       string
	OldValue   interface{}
	NewValue   interface{}
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
			Name:       dequote(strings.TrimSpace(attribute[0])),
			OldValue:   nil,
			NewValue:   doTypeConversion(attribute[1]),
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)
		if len(attribute) == 1 {
			// line does not have an "="
			// assume delimited with a space
			attribute = strings.Split(attribute[0], " ")
			return &AttributeChange{
				Name:       dequote(attribute[0]),
				OldValue:   doTypeConversion(attribute[len(attribute)-1]),
				NewValue:   nil,
				UpdateType: DestroyResource,
			}, nil
		}

		values := strings.Split(attribute[1], ATTRIBUTE_CHANGE_DELIMITER)
		if len(values) != 2 {
			return &AttributeChange{
				Name:       dequote(strings.TrimSpace(attribute[0])),
				OldValue:   doTypeConversion(strings.TrimSpace(attribute[1])),
				NewValue:   nil,
				UpdateType: DestroyResource,
			}, nil
		}

		return &AttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			OldValue:   doTypeConversion(values[0]),
			NewValue:   nil,
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
			if values[0] != SENSITIVE_VALUE {
				return nil, fmt.Errorf("failed to read attribute change from line %s", line)
			}
			values = append(values, SENSITIVE_VALUE)
		}

		return &AttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			OldValue:   doTypeConversion(values[0]),
			NewValue:   doTypeConversion(values[1]),
			UpdateType: updateType,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern %s", line)
	}
}

// NewAttributeChangeFromLine initializes an AttributeChange from a line within an Array attribute
// In an array resource, the attribute change does not have a name
func NewAttributeChangeFromArray(line string) (*AttributeChange, error) {
	line = strings.TrimSpace(line)
	if line == "" || line == "}" || IsResourceChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize an attributeChange", line)
	}
	if strings.HasPrefix(line, "+") {
		// add
		return &AttributeChange{
			OldValue:   nil,
			NewValue:   normalizeArrayAttribute(line),
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		return &AttributeChange{
			OldValue:   normalizeArrayAttribute(line),
			NewValue:   nil,
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		// TODO: confirm this is possible? I think array entries are immutable
		return nil, fmt.Errorf("unexpected replace single attribute in array %s", line)
	} else {
		return &AttributeChange{
			OldValue:   normalizeArrayAttribute(line),
			NewValue:   normalizeArrayAttribute(line),
			UpdateType: NoOpResource,
		}, nil
	}
}

// IsSensitive returns true if the attribute contains a sensitive value
func (a *AttributeChange) IsSensitive() bool {
	return a.OldValue == SENSITIVE_VALUE || a.NewValue == SENSITIVE_VALUE
}

// IsComputed returns true if the attribute contains a computed value
func (a *AttributeChange) IsComputed() bool {
	return a.OldValue == COMPUTED_VALUE || a.NewValue == COMPUTED_VALUE
}

func doTypeConversion(input string) interface{} {
	// if it has quotes, assume it is a string and return it without quotes
	if strings.HasPrefix(input, `"`) && strings.HasSuffix(input, `"`) {
		return dequote(input)
	}

	if input == "{}" {
		return nil
	}

	if input == "true" || input == "false" {
		b, err := strconv.ParseBool(input)
		if err != nil {
			return input
		}
		return b
	}

	if i, err := strconv.Atoi(input); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(input, 64); err == nil {
		return f
	}

	return input
}

func normalizeArrayAttribute(line string) interface{} {
	return doTypeConversion(strings.TrimRight(removeChangeTypeCharacters(line), ","))
}

func removeChangeTypeCharacters(line string) string {
	return strings.TrimLeft(line, "+/-~<= ")
}

func dequote(line string) string {
	return strings.Trim(line, "\"")
}
