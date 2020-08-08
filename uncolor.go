package tfplanparse

import (
	"bytes"

	"github.com/mattn/go-colorable"
)

func uncolor(in []byte) string {
	var out bytes.Buffer
	uncolorize := colorable.NewNonColorable(&out)
	uncolorize.Write(in)

	return out.String()
}
