package main

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/lipgloss"
	"github.com/sojamann/timg"
	"gopkg.in/yaml.v3"
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

func commentHandler(arg string, code string, maxWidth int) string {
	return ""
}

func imgHandler(arg string, block string, maxWidth int) string {
	fp, err := os.Open(abspath(strings.TrimSpace(block)))

	if err != nil {
		return lipgloss.NewStyle().
			MaxWidth(maxWidth).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).
			Render("IMAGE NOT FOUND")
	}
	defer fp.Close()

	img, _, err := image.Decode(fp)
	if err != nil {
		die(err.Error())
	}

	type imgConfig struct {
		Width  int
		Height int
	}

	var conf imgConfig
	if err := yaml.Unmarshal([]byte(arg), &conf); err != nil {
		conf.Width = maxWidth
		conf.Height = maxWidth / 2
	}

	return timg.Render(img, timg.FitTo(conf.Width, conf.Height))
}

// a blockHandler is can have some special logic to handle
// a certain block of code.
// args:
//
//	blockHandlerArgument - something to customize the behavior
//	section				 - the text/block to handle/render
//	maxWidth			 - the max width the block is allowed to have
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
	"img":     imgHandler,
	"comment": commentHandler,
}
