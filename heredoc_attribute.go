package tfplanparse

import (
	"fmt"
	"strings"
)

type HeredocAttributeChange struct {
	Name       string
	Token      string
	Before     []string
	After      []string
	UpdateType UpdateType
}

// IsHeredocAttributeChangeLine returns true if the line is a valid attribute change
// This requires the line to start with "+", "-" or "~", delimited with a space, and the value to start with "<<".
func IsHeredocAttributeChangeLine(line string) bool {
	line = strings.TrimSpace(line)
	attribute := strings.SplitN(line, ATTRIBUTE_DEFINITON_DELIMITER, 2)
	if len(attribute) != 2 {
		return false
	}

	validPrefix := strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "~")
	// The only permitted heredoc string is <<~EOT
	// Ref: https://github.com/hashicorp/terraform/blob/6a126df0c601ab23689171506bfc1386fea4c96c/command/format/diff.go#L838
	isHeredoc := strings.HasPrefix(attribute[1], "<<~EOT")
	return validPrefix && isHeredoc && !IsResourceChangeLine(line)
}

// IsHeredocAttributeTerminator returns true if the line is "EOT"
// EOT is the only possible terminator for the heredoc
// Ref: https://github.com/hashicorp/terraform/blob/6a126df0c601ab23689171506bfc1386fea4c96c/command/format/diff.go#L880
func IsHeredocAttributeTerminator(line string) bool {
	return strings.TrimSuffix(strings.TrimSpace(line), " -> null") == "EOT"
}

// NewHeredocAttributeChangeFromLine initializes a HeredocAttributeChange from a line containing a heredoc change
// It expects a line that passes the IsHeredocAttributeChangeLine check
func NewHeredocAttributeChangeFromLine(line string) (*HeredocAttributeChange, error) {
	line = strings.TrimSpace(line)
	if !IsHeredocAttributeChangeLine(line) {
		return nil, fmt.Errorf("%s is not a valid line to initialize a HeredocAttributeChange", line)
	}
	attribute := strings.SplitN(removeChangeTypeCharacters(line), ATTRIBUTE_DEFINITON_DELIMITER, 2)

	if strings.HasPrefix(line, "+") {
		// add
		return &HeredocAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			Before:     []string{},
			After:      []string{},
			UpdateType: NewResource,
		}, nil
	} else if strings.HasPrefix(line, "-") {
		// destroy
		return &HeredocAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			Before:     []string{},
			After:      []string{},
			UpdateType: DestroyResource,
		}, nil
	} else if strings.HasPrefix(line, "~") {
		// replace
		updateType := UpdateInPlaceResource
		if strings.HasSuffix(attribute[1], " # forces replacement") {
			updateType = ForceReplaceResource
		}

		return &HeredocAttributeChange{
			Name:       dequote(strings.TrimSpace(attribute[0])),
			Before:     []string{},
			After:      []string{},
			UpdateType: updateType,
		}, nil
	} else {
		return nil, fmt.Errorf("unrecognized line pattern")
	}
}

func (h *HeredocAttributeChange) AddLineToContent(line string) {
	switch h.UpdateType {
	case NewResource:
		h.After = append(h.After, strings.TrimPrefix(line, "+ "))
	case DestroyResource:
		h.Before = append(h.Before, strings.TrimPrefix(line, "- "))
	default:
		// replace
		// TODO: remove ANSI coded - or + only
		h.Before = append(h.Before, line)
		h.After = append(h.After, line)
	}
}

func (h *HeredocAttributeChange) GetBeforeAttribute(opts ...GetBeforeAfterOptions) string {
	return strings.Join(h.Before, "\n")
}

func (h *HeredocAttributeChange) GetAfterAttribute(opts ...GetBeforeAfterOptions) string {
	return strings.Join(h.After, "\n")
}
