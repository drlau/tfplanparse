package tfplanparse

import (
	"github.com/drlau/tf-plan-parse/plan"
)

type PlanResult struct {
	Resources []*plan.ResourcePlan
}