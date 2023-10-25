package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}

// TODO: rename
func codeHandler(arg string, code string) string {
	lexer := lexers.Fallback
	// try getting the user defined lexer based on the name of the
	// language. If not use the fallback lexer instead
	if choice := lexers.Get(arg); choice != nil {
		lexer = choice
	}
	// the code starts at the next line
	buff := &strings.Builder{}
	quick.Highlight(buff, code, lexer.Config().Name, "terminal", "monokai")
	return buff.String()
}

var pluginRegEx = regexp.MustCompile("@(code)(:{(.*)})?\n")
var handlers = map[string]func(string, string) string{
	"code": codeHandler,
}

type model struct {
	slides []string

	currSlide int

	ready bool
	vp    viewport.Model
}

func (m *model) Init() tea.Cmd {
	return nil
}

// renders normal text
var namedStyles = map[string]lipgloss.Style{
	"h1": lipgloss.NewStyle().
		Bold(true).
		Background(lipgloss.Color("#F25D94")),
	"b": lipgloss.NewStyle().Bold(true),
}
var namedStyleNames = strings.Join(Keys(namedStyles), "|")
var h1Regex = regexp.MustCompile(fmt.Sprintf("!(%s){(.+?)}", namedStyleNames))

func (m *model) renderText(text string) string {

	buff := &strings.Builder{}

	var offset int
	for {
		indecies := h1Regex.FindStringSubmatchIndex(text[offset:])
		// no plugin found .... render the text
		if len(indecies) == 0 {
			buff.WriteString(text)
			break
		}

		// write out everything we have seen up till the style starts
		buff.WriteString(text[:indecies[0]])

		stylename := text[indecies[2] : indecies[3]]
		toStyle := text[indecies[4] : indecies[5]]

		buff.WriteString(namedStyles[stylename].Render(toStyle))
		text=text[indecies[1]:]
	}
	return buff.String()
}

// TODO: add error handling
func (m *model) renderSlide(slide string) string {
	buff := &strings.Builder{}

	var offset int
	for {
		indecies := pluginRegEx.FindStringSubmatchIndex(slide[offset:])
		// no plugin found .... render the text
		if len(indecies) == 0 {
			buff.WriteString(m.renderText(slide[offset:]))
			break
		}

		// write out everything we have seen up till the plugin starts
		beforeString := m.renderText(slide[offset : offset+indecies[0]])
		buff.WriteString(beforeString)

		pluginName := slide[offset+indecies[2] : offset+indecies[3]]
		pluginOpt := ""
		if indecies[4] != -1 {
			pluginOpt = slide[offset+indecies[6] : offset+indecies[7]]
		}

		// lets get the content
		offset += indecies[1]
		contentStart := strings.Index(slide[offset:], "```") + offset + 3
		offset = contentStart

		contentEnd := strings.Index(slide[offset:], "```") + offset
		offset = contentEnd + 3

		pluginArg := slide[contentStart:contentEnd]
		handlerResult := handlers[pluginName](pluginOpt, pluginArg)
		buff.WriteString(handlerResult)
	}

	// we want to center the entire content box not the
	// text itself (this would be alignment)
	width, height := lipgloss.Size(buff.String())
	return lipgloss.NewStyle().
		Height(m.vp.Height - 1).
		MaxHeight(m.vp.Height - 1).
		Width(m.vp.Width).
		MaxWidth(m.vp.Width).
		MarginTop((m.vp.Height - height) / 2).
		MarginBottom((m.vp.Height - height) / 2).
		MarginLeft((m.vp.Width - width) / 2).
		MarginRight((m.vp.Width - width) / 2).
		Render(buff.String())
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "q":
			return m, tea.Quit
		case "l", tea.KeyRight.String():
			m.currSlide = min(m.currSlide+1, len(m.slides)-1)
			m.vp.SetContent(m.renderSlide(m.slides[m.currSlide]))
		case "h", tea.KeyLeft.String():
			m.currSlide = max(0, m.currSlide-1)
			m.vp.SetContent(m.renderSlide(m.slides[m.currSlide]))
		}

	case tea.WindowSizeMsg:
		// -1 because of footer
		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-1)
			m.vp.SetContent(m.renderSlide(m.slides[m.currSlide]))
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 1
		}
		m.ready = true
	}

	// THe viewport might wants to do some things
	m.vp, cmd = m.vp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	if !m.ready {
		return "\n Loading ... "
	}

	parts := 2
	progress := lipgloss.NewStyle().
		Width(m.vp.Width / parts).
		Align(lipgloss.Left).
		PaddingLeft(2).
		Render(fmt.Sprintf("[%d/%d]", m.currSlide+1, len(m.slides)))

	author := lipgloss.NewStyle().
		Width(m.vp.Width / parts).
		Align(lipgloss.Right).
		PaddingRight(2).
		Render("Author Name")
	status := progress + author

	return fmt.Sprintf("%s\n%s", m.vp.View(), status)
}

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	slides := strings.Split(string(data), "---")
	if len(slides) == 0 {
		return
	}
	if slides[0] == "" {
		slides = slides[1:]
	}
	m := &model{
		slides: slides,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
