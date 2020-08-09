package tfplanparse

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const (
	NO_CHANGES_STRING    = "No changes. Infrastructure is up-to-date."
	CHANGES_START_STRING = "Terraform will perform the following actions:"
	CHANGES_END_STRING   = "Plan: "
	ERROR_STRING         = "Error: "
)

func Parse(input io.Reader) ([]*ResourceChange, error) {
	result := []*ResourceChange{}
	var resourceChange *ResourceChange
	var mapAttributeChange *MapAttributeChange
	var err error

	parse := false
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		text := strings.TrimSpace(uncolor(scanner.Bytes()))
		if text == "" {
			continue
		}

		if !parse {
			if strings.Contains(text, NO_CHANGES_STRING) || strings.Contains(text, ERROR_STRING) {
				// Nothing to parse, return empty plan
				return result, nil
			} else if strings.Contains(text, CHANGES_START_STRING) {
				// Parse all lines from here on
				parse = true
			}

			continue
		}

		if strings.Contains(text, CHANGES_END_STRING) {
			// we are done
			if resourceChange != nil {
				result = append(result, resourceChange)
			}

			return result, nil
		}

		if IsResourceCommentLine(text) {
			// if parsing a resource before, append it as is
			if resourceChange != nil {
				result = append(result, resourceChange)
			}

			resourceChange, err = NewResourceChangeFromComment(text)
			if err != nil {
				return result, err
			}
			// TODO: handle nested maps
		} else if IsMapAttributeChangeLine(text) {
			mapAttributeChange, err = NewMapAttributeChangeFromLine(text)
			if err != nil {
				return result, err
			}
		} else if IsAttributeChangeLine(text) {
			ac, err := NewAttributeChangeFromLine(text)
			if err != nil {
				return result, err
			}

			// if currently parsing a map attribute, this attribute belongs to the map
			if mapAttributeChange != nil {
				mapAttributeChange.AttributeChanges = append(mapAttributeChange.AttributeChanges, ac)
			} else {
				resourceChange.AttributeChanges = append(resourceChange.AttributeChanges, ac)
			}
			// TODO: this does not handle nested maps at all
		} else if mapAttributeChange != nil && IsMapAttributeTerminator(text) {
			if resourceChange != nil {
				resourceChange.MapAttributeChanges = append(resourceChange.MapAttributeChanges, mapAttributeChange)
				mapAttributeChange = nil
			}
		}
	}

	if resourceChange != nil {
		result = append(result, resourceChange)
	}

	return result, nil
}

func ParseFromFile(filepath string) ([]*ResourceChange, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return []*ResourceChange{}, err
	}

	return Parse(f)
}
