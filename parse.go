package tfplanparse

import (
	"bufio"
	"fmt"
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
	parse := false
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		text := formatInput(scanner.Bytes())
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

		if IsResourceCommentLine(text) {
			rc, err := parseResource(scanner)
			if err != nil {
				return nil, err
			}

			result = append(result, rc)
		}

		if strings.Contains(formatInput(scanner.Bytes()), CHANGES_END_STRING) {
			// we are done
			return result, nil
		}
	}

	return nil, fmt.Errorf("unexpected end of input while parsing plan")
}

func ParseFromFile(filepath string) ([]*ResourceChange, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return []*ResourceChange{}, err
	}

	return Parse(f)
}

func parseResource(s *bufio.Scanner) (*ResourceChange, error) {
	rc, err := NewResourceChangeFromComment(s.Text())
	if err != nil {
		return nil, err
	}
	for s.Scan() {
		text := formatInput(s.Bytes())
		switch {
		case IsResourceCommentLine(text), strings.Contains(text, CHANGES_END_STRING):
			return rc, nil
		case IsMapAttributeChangeLine(text):
			ma, err := parseMapAttribute(s)
			if err != nil {
				return nil, err
			}
			rc.MapAttributeChanges = append(rc.MapAttributeChanges, ma)
		case IsAttributeChangeLine(text):
			ac, err := NewAttributeChangeFromLine(text)
			if err != nil {
				return nil, err
			}
			rc.AttributeChanges = append(rc.AttributeChanges, ac)
		}
	}

	return nil, fmt.Errorf("unexpected end of input while parsing resource")
}

func parseMapAttribute(s *bufio.Scanner) (*MapAttributeChange, error) {
	result, err := NewMapAttributeChangeFromLine(s.Text())
	if err != nil {
		return nil, err
	}
	for s.Scan() {
		text := formatInput(s.Bytes())
		switch {
		case IsMapAttributeTerminator(text):
			return result, nil
		case IsResourceCommentLine(text), strings.Contains(text, CHANGES_END_STRING):
			return nil, fmt.Errorf("unexpected line while parsing map attribute: %s", text)
		case IsMapAttributeChangeLine(text):
			ma, err := parseMapAttribute(s)
			if err != nil {
				return nil, err
			}
			result.MapAttributeChanges = append(result.MapAttributeChanges, ma)
		case IsAttributeChangeLine(text):
			ac, err := NewAttributeChangeFromLine(text)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ac)
		}
	}

	return nil, fmt.Errorf("unexpected end of input while parsing map attribute")
}

func formatInput(input []byte) string {
	return strings.TrimSpace(uncolor(input))
}