package main

import (
	"github.com/charmbracelet/lipgloss"
)

// renders normal text
var namedStyleLookupTable = map[string]lipgloss.Style{
	"h1": lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#F25D94")),
	"b": lipgloss.NewStyle().Bold(true),
	"i": lipgloss.NewStyle().Italic(true),
}
