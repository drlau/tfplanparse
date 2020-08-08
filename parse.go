package tfplanparse

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

const (
	NO_CHANGES_STRING    = "No changes. Infrastructure is up-to-date."
	CHANGES_START_STRING = "Terraform will perform the following actions:"
	CHANGES_END_STRING   = "Plan: "
	ERROR_STRING         = "Error: "
)

func Parse(input io.Reader) []*ResourcePlan {
	// Overall:
	// Look for the start of resources
	// No changes -> return
	// Changes -> start parsing
	// New / Force Replace -> Parse every line
	// Update in place -> parse only changed lines
	// Destroy -> name only
	return []*ResourcePlan{}
}

// TODO: handle multi level structs
func ParseFromFile(filepath string) []*ResourcePlan {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	result := []*ResourcePlan{}
	var resourcePlan *ResourcePlan
	var mapAttributeChange *MapAttributeChange

	parse := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(uncolor(scanner.Bytes()))

		if parse {
			if text == "" {
				continue
			} else if strings.Contains(text, CHANGES_END_STRING) {
				// we are done

				if resourcePlan != nil {
					result = append(result, resourcePlan)
				}
				return result
			}

			if IsResourceCommentLine(text) {
				if resourcePlan != nil {
					result = append(result, resourcePlan)
				}

				resourcePlan, err = NewResourcePlanFromComment(text)
				if err != nil {
					panic(err)
				}
			} else if IsMapAttributeChangeLine(text) {
				mapAttributeChange, err = NewMapAttributeChangeFromLine(text)
				if err != nil {
					panic(err)
				}
			} else if IsAttributeChangeLine(text) {
				log.Printf("running for line %v\n", text)
				ac, err := NewAttributeChangeFromLine(text)
				if err != nil {
					panic(err)
				}
				if mapAttributeChange != nil {
					mapAttributeChange.AttributeChanges = append(mapAttributeChange.AttributeChanges, ac)
				} else {
					resourcePlan.AttributeChanges = append(resourcePlan.AttributeChanges, ac)
				}
			} else if mapAttributeChange != nil && IsMapAttributeTerminator(text) {
				if resourcePlan != nil {
					resourcePlan.MapAttributeChanges = append(resourcePlan.MapAttributeChanges, mapAttributeChange)
					mapAttributeChange = nil
				}
			} else {
				log.Printf("skipping line: %s\n", text)
			}
		} else {
			if strings.Contains(text, NO_CHANGES_STRING) || strings.Contains(text, ERROR_STRING) {
				// Nothing to parse, return empty plan
				return result
			} else if strings.Contains(text, CHANGES_START_STRING) {
				// Parse all lines from here on
				parse = true
			}
		}
	}

	if resourcePlan != nil {
		result = append(result, resourcePlan)
	}

	return result
}
