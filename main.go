package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/keymaps"
)

func main() {
	initKeys()
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

func initKeys() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	configDir := filepath.Join(homeDir, ".config", "helm-tui")
	configPath := filepath.Join(configDir, "config.json")

	// 1. Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}

	// 2. Check if there is a config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return
	}

	// 3. Read + parse safely
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	tmp := keymaps.DefaultConfig // start with defaults

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tmp); err != nil {
		fmt.Println("Error parsing config, using defaults:", err)
		return
	}

	fmt.Println(keymaps.DefaultConfig)

	// 4. Apply config
	keymaps.DefaultConfig = tmp
}
