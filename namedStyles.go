package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renders normal text
var namedStyleLookupTable = map[string]lipgloss.Style{
	"h1": lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#F25D94")),
	"b": lipgloss.NewStyle().Bold(true),
}
var namedStyleNames = strings.Join(Keys(namedStyleLookupTable), "|")
var h1Regex = regexp.MustCompile(fmt.Sprintf("!(%s){(.+?)}", namedStyleNames))

