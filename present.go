package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	slides []string
	currSlide int

	namedStyles map[string]lipgloss.Style
	blockHandlers map[string]blockHandler

	author string

	ready bool
	vp    viewport.Model
}

func (m *model) Init() tea.Cmd {
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

		stylename := text[indecies[2] : indecies[3]]
		toStyle := text[indecies[4] : indecies[5]]

		style, found := m.namedStyles[stylename];
		if !found {
			// NOTE: bubbletea does not like this...
			log.Fatalln(fmt.Errorf("You used the style %s which is neither builtin or custom", stylename))
		}

		buff.WriteString(style.MaxWidth(width).Render(toStyle))

		text=text[indecies[1]:]
	}
	return buff.String()
}


var blockHandlerNamePat = `(.+?)`
var blockHandlerOptPat = `(\[(.*?)\])?`
var blockHandlerContentPat = `{((?s).*?)}` // (?s) makes . match \n
var blockHandlerRegex= regexp.MustCompile(fmt.Sprintf("@%s%s%s", blockHandlerNamePat, blockHandlerOptPat, blockHandlerContentPat))

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
			log.Fatalln(fmt.Errorf("block handler not defined"))	
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

