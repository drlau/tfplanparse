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
	rc, err := NewResourceChangeFromComment(formatInput(s.Bytes()))
	if err != nil {
		return nil, err
	}
	for s.Scan() {
		text := formatInput(s.Bytes())
		switch {
		case IsResourceTerminator(text):
			return rc, nil
		case IsResourceCommentLine(text), strings.Contains(text, CHANGES_END_STRING):
			return nil, fmt.Errorf("unexpected line while parsing resource attribute: %s", text)
		case IsMapAttributeChangeLine(text):
			ma, err := parseMapAttribute(s)
			if err != nil {
				return nil, err
			}
			rc.AttributeChanges = append(rc.AttributeChanges, ma)
		case IsArrayAttributeChangeLine(text):
			aa, err := parseArrayAttribute(s)
			if err != nil {
				return nil, err
			}
			rc.AttributeChanges = append(rc.AttributeChanges, aa)
		case IsHeredocAttributeChangeLine(text):
			ha, err := parseHeredocAttribute(s)
			if err != nil {
				return nil, err
			}
			rc.AttributeChanges = append(rc.AttributeChanges, ha)
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
	normalized := formatInput(s.Bytes())
	result, err := NewMapAttributeChangeFromLine(normalized)
	if err != nil {
		return nil, err
	}
	if IsOneLineEmptyMapAttribute(normalized) {
		return result, nil
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
			result.AttributeChanges = append(result.AttributeChanges, ma)
		case IsArrayAttributeChangeLine(text):
			aa, err := parseArrayAttribute(s)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, aa)
		case IsHeredocAttributeChangeLine(text):
			ha, err := parseHeredocAttribute(s)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ha)
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

func parseArrayAttribute(s *bufio.Scanner) (*ArrayAttributeChange, error) {
	normalized := formatInput(s.Bytes())
	result, err := NewArrayAttributeChangeFromLine(normalized)
	if err != nil {
		return nil, err
	}
	if IsOneLineEmptyArrayAttribute(normalized) {
		return result, nil
	}
	// TODO: all elements of array attributes are the same type
	for s.Scan() {
		text := formatInput(s.Bytes())
		switch {
		case IsArrayAttributeTerminator(text):
			return result, nil
		case IsResourceCommentLine(text), strings.Contains(text, CHANGES_END_STRING):
			return nil, fmt.Errorf("unexpected line while parsing array attribute: %s", text)
		case IsMapAttributeChangeLine(text):
			ma, err := parseMapAttribute(s)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ma)
		case IsArrayAttributeChangeLine(text):
			ma, err := parseArrayAttribute(s)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ma)
		case IsHeredocAttributeChangeLine(text):
			ha, err := parseHeredocAttribute(s)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ha)
		case IsAttributeChangeArrayItem(text):
			ac, err := NewAttributeChangeFromArray(text)
			if err != nil {
				return nil, err
			}
			result.AttributeChanges = append(result.AttributeChanges, ac)
		}
	}

	return nil, fmt.Errorf("unexpected end of input while parsing array attribute")
}

func parseHeredocAttribute(s *bufio.Scanner) (*HeredocAttributeChange, error) {
	normalized := formatInput(s.Bytes())
	result, err := NewHeredocAttributeChangeFromLine(normalized)
	if err != nil {
		return nil, err
	}
	for s.Scan() {
		text := formatInput(s.Bytes())
		if IsHeredocAttributeTerminator(text) {
			return result, nil
		}

		// TODO: should not trim space for heredoc, but only trim the indent
		// TODO: it's also hard to determine if a line was deleted or is a list in an array
		result.AddLineToContent(text)
	}

	return nil, fmt.Errorf("unexpected end of input while parsing heredoc attribute")
}

func formatInput(input []byte) string {
	return strings.TrimSpace(uncolor(input))
}
