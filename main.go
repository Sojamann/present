package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
			m.vp.SetContent(m.slides[m.currSlide])
		case "h", tea.KeyLeft.String():
			m.currSlide = max(0, m.currSlide-1)
			m.vp.SetContent(m.slides[m.currSlide])
		}

	case tea.WindowSizeMsg:
		// -1 because of footer
		if !m.ready { 
			m.vp = viewport.New(msg.Width, msg.Height-1)
			m.vp.SetContent(m.slides[m.currSlide])
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height-1
		}
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

	footer := lipgloss.NewStyle().Render("q - Quit | h/Left - prev slide | l/Right - next slide")

	return fmt.Sprintf("%s\n%s", m.vp.View(), footer)
}


func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	m := &model {
		slides: strings.Split(string(data), "---"),
	}
	p := tea.NewProgram(m, tea.WithAltScreen());
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
