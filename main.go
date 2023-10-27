package main

import (
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)


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

	m := NewPresentation(config, slides)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
