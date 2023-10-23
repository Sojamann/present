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

func (m *model) renderSlide(slide string) string {
	return lipgloss.NewStyle().
		Width(m.vp.Width-1).
		MaxWidth(m.vp.Width-1).
		Height(m.vp.Height-1).
		MaxHeight(m.vp.Height-1).
		Align(lipgloss.Center).
		Margin(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		BorderRight(true).
		Render(slide)
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
