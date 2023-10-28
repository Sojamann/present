package main

import (
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/lipgloss"
)

// TODO: rename
func codeHandler(arg string, code string, width int) string {
	lexer := lexers.Fallback
	// try getting the user defined lexer based on the name of the
	// language. If not use the fallback lexer instead
	if choice := lexers.Get(arg); choice != nil {
		lexer = choice
	}
	// the code starts at the next line
	buff := &strings.Builder{}
	quick.Highlight(buff, code, lexer.Config().Name, "terminal", "monokai")
	return lipgloss.NewStyle().MaxWidth(width).Render(buff.String())
}

func noteHandler(arg string, code string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"}).
		Padding(1).
		Render(code)
}

func warningHandler(arg string, code string, width int) string {
	return lipgloss.NewStyle().
		Width(width).
		MaxWidth(width).
		Foreground(lipgloss.Color("#FFFDF5")).
		Background(lipgloss.Color("#FF5F87")).
		Padding(1).
		Render(code)

}

// a blockHandler is can have some special logic to handle
// a certain block of code.
// args:
//
//	blockHandlerArgument - something to customize the behavior
//	section				 - the text/block to handle/render
//	width				 - the max width the block is allowed to have
//							knowing that the blockHandler is able to make
//							some better adjustments than when MaxWidth is
//							applied from the outside
//
// returns:
//
//	the string which should be displyed
type blockHandler func(string, string, int) string

var DefaultBlockHandlers = map[string]blockHandler{
	"code":    codeHandler,
	"note":    noteHandler,
	"warning": warningHandler,
}
