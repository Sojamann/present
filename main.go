package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
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
	configSlide, slides := slides[0], slides[1:]
	var config PresentationConfig
	yaml.Unmarshal([]byte(configSlide), &config)
	
	//fmt.Println(config)
	//os.Exit(0)

	customStyles := customStyleToLipgloss(config.Style)

	mergeStyles := MapMerge(customStyles, namedStyleLookupTable)
	author := config.Author

	m := &model{
		slides: slides,
		author: author,
		namedStyles: mergeStyles,
	}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
