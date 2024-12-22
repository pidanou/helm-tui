package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/helpers"
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	helpers.LogFile = f
	defer f.Close()
	defer os.Truncate("debug.log", 0)

	var tabs []string

	// Iterate over the map and collect the values
	for _, value := range tabLabels {
		tabs = append(tabs, value)
	}

	p := tea.NewProgram(newModel(tabs), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
