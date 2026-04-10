package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	initKeys()

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

func initKeys() {
	// Get the user's home directory (cross-platform)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	// Build the config file path
	configPath := filepath.Join(homeDir, ".config", "helm-tui", "config")

	// Open the file
	file, err := os.Open(configPath)
	if err != nil {
		return
	}
	defer file.Close()

	// Parse JSON
	var config map[string]any
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	// Print parsed result
	fmt.Printf("Parsed config: %+v\n", config)
}
