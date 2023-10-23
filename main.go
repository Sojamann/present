package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	slides []string

	currSlide int

	ready bool
	vp viewport.Model
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) renderText(text string) string {
	return lipgloss.NewStyle().
		Width(m.vp.Width).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(text)
}

func (m *model) renderCode(code string) string {
	start := strings.Index(code, "\n")
	lexer := lexers.Fallback
	// try getting the user defined lexer based on the name of the 
	// language. If not use the fallback lexer instead
	if choice := lexers.Get(code[:start]); choice != nil {
		lexer = choice
	}
	// the code starts at the next line
	code = code[start:]
	buff := &strings.Builder{}
	quick.Highlight(buff, code, lexer.Config().Name, "terminal", "monokai")
	return lipgloss.NewStyle().
		AlignVertical(lipgloss.Center).
		Render(buff.String())
}

func (m *model) renderSlide(slide string) string {
	buff := &strings.Builder{}
	
	var last, found int
	var insideCode bool
outer:
	for {
		found = strings.Index(slide[last:], "```")
		if found == -1 {
			buff.WriteString(m.renderText(slide[last:]))
			break outer
		}

		found += last

		if insideCode {
			buff.WriteString(m.renderCode(slide[last:found]))
			insideCode = false
		}else { 
			buff.WriteString(m.renderText(slide[last:found]))	
			insideCode = true
		}
		last = found+3 // go ``` many ahead
	}
	
	return lipgloss.NewStyle().
		Height(m.vp.Height-1).
		MaxHeight(m.vp.Height-1).
		Padding(1).
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
			println("RENDER SLIDE", m.renderSlide(m.slides[m.currSlide]))
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height-1
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
		Width(m.vp.Width/parts).
		Align(lipgloss.Left).
		PaddingLeft(2).
		Render(fmt.Sprintf("[%d/%d]", m.currSlide+1, len(m.slides)))

	author := lipgloss.NewStyle().
		Width(m.vp.Width/parts).
		Align(lipgloss.Right).
		PaddingRight(2).
		Render("Author Name")
	status := progress+author 

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
	m := &model {
		slides: slides,
	}
	p := tea.NewProgram(m, tea.WithAltScreen());
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
