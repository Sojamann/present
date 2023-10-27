package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type styleConfig struct {
	Bold bool
	Italic bool
	Fg string
	Bg string
}

type PresentationConfig struct {
	Author string
	Style map[string]styleConfig
}

func customStyleToLipgloss(custom map[string]styleConfig) map[string]lipgloss.Style {
	result := make(map[string]lipgloss.Style)

	for name, config := range custom {
		if _, found := namedStyleLookupTable[name]; found {
			log.Fatalln(fmt.Errorf("The custom name %s is a named style", name))
		}

		result[name] = lipgloss.NewStyle().
			Bold(config.Bold).
			Italic(config.Italic).
			Foreground(lipgloss.Color(config.Fg)).
			Background(lipgloss.Color(config.Bg))
	}
	return result
}

func NewPresentation(config PresentationConfig, slides []string) *model {
	return &model{
		slides: slides,
		author: config.Author,
		customStyle: customStyleToLipgloss(config.Style),
	}	
}

type model struct {
	slides []string
	config PresentationConfig
	currSlide int

	customStyle map[string]lipgloss.Style

	author string

	ready bool
	vp    viewport.Model
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) renderText(text string, width int) string {
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
		
		if style, found := m.customStyle[stylename]; found {
			buff.WriteString(style.MaxWidth(width).Render(toStyle))
		} else if style, found := namedStyleLookupTable[stylename]; found {
			buff.WriteString(style.MaxWidth(width).Render(toStyle))
		} else {
			// NOTE: bubbletea does not like this...
			log.Fatalln(fmt.Errorf("You used the style %s which is neither builtin or custom", stylename))
		}

		text=text[indecies[1]:]
	}
	return buff.String()
}

// TODO: add error handling
func (m *model) renderSlide(slide string) string {
	buff := &strings.Builder{}

	var offset int
	for {
		indecies := blockHandlerRegex.FindStringSubmatchIndex(slide[offset:])
		// no plugin found .... render the text
		if len(indecies) == 0 {
			buff.WriteString(m.renderText(slide[offset:], m.vp.Width))
			break
		}

		// write out everything we have seen up till the plugin starts
		beforeString := m.renderText(slide[offset : offset+indecies[0]], m.vp.Width)
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
		handlerResult := blockHandlerLookupTable[pluginName](pluginOpt, pluginArg, m.vp.Width)
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
		case "q", tea.KeyEscape.String():
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

