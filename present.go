package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	slides    []string
	currSlide int

	namedStyles   map[string]lipgloss.Style
	blockHandlers map[string]blockHandler

	author string

	ready bool
	vp    viewport.Model

	errors []string
}

func (m *model) Init() tea.Cmd {
	// validate slides
	for i, slide := range m.slides {
		for _, match := range namedStyleRegex.FindAllStringSubmatch(slide, -1) {
			if _, found := m.namedStyles[match[1]]; !found {
				m.errors = append(m.errors, fmt.Sprintf("Slide %d - unknown style '%s' in: %s", i, match[1], match[0])) 
			}
		}
		for _, match := range blockHandlerRegex.FindAllStringSubmatch(slide, -1) {
			if _, found := m.blockHandlers[match[1]]; !found {
				m.errors = append(m.errors, fmt.Sprintf("Slide %d - unknown block handler '%s' in: %s", i, match[1], match[0])) 
			}
		}
	}

	return nil
}

var namedStyleNamePat = `(.+?)`
var namedStyleContentPat = `((?s).+?)`
var namedStyleRegex = regexp.MustCompile(fmt.Sprintf("!%s{%s}", namedStyleNamePat, namedStyleContentPat))

func (m *model) renderText(text string, width int) string {
	buff := &strings.Builder{}

	for {
		// 0/1  - start/end
		// 2/3	- style name
		// 4/5	- content
		indecies := namedStyleRegex.FindStringSubmatchIndex(text)
		// no plugin found .... render the text
		if len(indecies) == 0 {
			buff.WriteString(text)
			break
		}

		// write out everything we have seen up till the style starts
		buff.WriteString(text[:indecies[0]])

		styleName := text[indecies[2]:indecies[3]]
		content := text[indecies[4]:indecies[5]]

		style, found := m.namedStyles[styleName]
		if !found {
			panic(fmt.Sprintf("Style '%s' not defined", styleName))
		}

		buff.WriteString(style.MaxWidth(width).Render(content))

		text = text[indecies[1]:]
	}
	return buff.String()
}

var blockHandlerNamePat = `(.+?)`
var blockHandlerOptPat = `(\[(.*?)\])?`
var blockHandlerContentPat = `{((?s).*?)}` // (?s) makes . match \n
var blockHandlerRegex = regexp.MustCompile(fmt.Sprintf("@%s%s%s", blockHandlerNamePat, blockHandlerOptPat, blockHandlerContentPat))

// TODO: add error handling
func (m *model) renderSlide(slide string) string {
	buff := &strings.Builder{}

	for {
		// 0/1 regexStart/End
		// 2/3 handlerName
		// 4/5 are options defined at all?
		// 6/7 the handler options/argString
		// 8/9 the block to handle
		indecies := blockHandlerRegex.FindStringSubmatchIndex(slide)
		// no blocks found .... render the text
		if len(indecies) == 0 {
			buff.WriteString(m.renderText(slide, m.vp.Width))
			break
		}

		// write out everything before the block starts
		beforeString := m.renderText(slide[:indecies[0]], m.vp.Width)
		buff.WriteString(beforeString)

		blockHandlerName := slide[indecies[2]:indecies[3]]
		blockHandlerOpt := ""
		if indecies[4] != -1 {
			blockHandlerOpt = slide[indecies[6]:indecies[7]]
		}

		blockContent := slide[indecies[8]:indecies[9]]

		blockHandler, found := m.blockHandlers[blockHandlerName]
		if !found {
			panic(fmt.Sprintf("Block Handler '%s' not defined", blockHandlerName))
		}

		handlerResult := blockHandler(blockHandlerOpt, blockContent, m.vp.Width)
		buff.WriteString(handlerResult)

		slide = slide[indecies[1]:]
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
		case "h", tea.KeyLeft.String():
			m.currSlide = max(0, m.currSlide-1)
		}
		
		if len(m.errors) == 0 {
			m.vp.SetContent(m.renderSlide(m.slides[m.currSlide]))
		}

	case tea.WindowSizeMsg:
		// -1 because of footer
		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-1)
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 1
		}
		if len(m.errors) == 0 {
			m.vp.SetContent(m.renderSlide(m.slides[m.currSlide]))
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
	if len(m.errors) > 0 {
		return lipgloss.NewStyle().
			Width(m.vp.Width).
			MaxWidth(m.vp.Width).
			MaxHeight(m.vp.Height).
			Padding(3).
			Foreground(lipgloss.Color("1")).
			Render(strings.Join(m.errors, "\n\n"))
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
		Render(m.author)
	status := progress + author

	return fmt.Sprintf("%s\n%s", m.vp.View(), status)
}
