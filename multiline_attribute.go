package tfplanparse

import (
	"strings"
)

// TODO: interface

func getMultiLineAttributeName(line string) string {
	line = removeChangeTypeCharacters(line)
	// Multiline attributes may or may not have a name
	// If they do have a name, they are delimited with a '=' or a ' '
	attribute := strings.SplitN(line, ATTRIBUTE_DEFINITON_DELIMITER, 2)
	if len(attribute) == 2 {
		return dequote(strings.TrimSpace(attribute[0]))
	}

	attribute = strings.SplitN(line, " ", 2)
	if len(attribute) == 2 {
		return dequote(strings.TrimSpace(attribute[0]))
	}

	return ""
}
