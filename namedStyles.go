package main

import (
	"github.com/charmbracelet/lipgloss"
)

// renders normal text
var namedStyleLookupTable = map[string]lipgloss.Style{
	"h": lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#F25D94")),
	"b": lipgloss.NewStyle().Bold(true),
	"i": lipgloss.NewStyle().Italic(true),

	// colors
	"red": lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
	"green": lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
	"yellow": lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
	"blue": lipgloss.NewStyle().Foreground(lipgloss.Color("4")),
	"white": lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
}
