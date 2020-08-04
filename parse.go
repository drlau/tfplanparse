package vision

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/drlau/vision/plan"
)

const (
	NO_CHANGES_STRING    = "No changes. Infrastructure is up-to-date."
	CHANGES_START_STRING = "Terraform will perform the following actions:"
	CHANGES_END_STRING   = "Plan: "
	ERROR_STRING         = "Error: "
)

// TODO: remove ANSI color codes

func Parse(body string) *PlanResult {
	// Overall:
	// Look for the start of resources
	// No changes -> return
	// Changes -> start parsing
	// New / Force Replace -> Parse every line
	// Update in place -> parse only changed lines
	// Destroy -> name only
	return &PlanResult{}
}

// TODO: handle multi level structs
func ParseFromFile(filepath string) *PlanResult {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	result := &PlanResult{}
	var resourcePlan *plan.ResourcePlan
	var mapAttributeChange *plan.MapAttributeChange

	parse := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())

		if parse {
			if text == "" {
				continue
			} else if strings.Contains(text, CHANGES_END_STRING) {
				// we are done

				if resourcePlan != nil {
					result.Resources = append(result.Resources, resourcePlan)
				}
				return result
			}

			if plan.IsResourceCommentLine(text) {
				if resourcePlan != nil {
					result.Resources = append(result.Resources, resourcePlan)
				}

				resourcePlan, err = plan.NewResourcePlanFromComment(text)
				if err != nil {
					panic(err)
				}
			} else if plan.IsMapAttributeChangeLine(text) {
				mapAttributeChange, err = plan.NewMapAttributeChangeFromLine(text)
				if err != nil {
					panic(err)
				}
			} else if plan.IsAttributeChangeLine(text) {
				log.Printf("running for line %v\n", text)
				ac, err := plan.NewAttributeChangeFromLine(text)
				if err != nil {
					panic(err)
				}
				if mapAttributeChange != nil {
					mapAttributeChange.AttributeChanges = append(mapAttributeChange.AttributeChanges, ac)
				} else {
					resourcePlan.AttributeChanges = append(resourcePlan.AttributeChanges, ac)
				}
			} else if mapAttributeChange != nil && plan.IsMapAttributeTerminator(text) {
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
		result.Resources = append(result.Resources, resourcePlan)
	}

	return result
}
