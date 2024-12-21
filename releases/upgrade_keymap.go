package releases

import "github.com/charmbracelet/bubbles/key"

var upgradeKeys = keyMap{
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}
