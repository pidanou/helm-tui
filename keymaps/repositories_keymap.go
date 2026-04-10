package keymaps

import (
	"github.com/charmbracelet/bubbles/key"
)

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type RepositoriesKeyMap struct {
	Add                      key.Binding
	Delete                   key.Binding
	Refresh                  key.Binding
	Move                     key.Binding
	Up                       key.Binding
	Down                     key.Binding
	Left                     key.Binding
	Right                    key.Binding
	Update                   key.Binding
	Install                  key.Binding
	Select                   key.Binding
	Cancel                   key.Binding
	ShowPackageDefaultValues key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k RepositoriesKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.ShowPackageDefaultValues, k.Delete, k.Update, k.Move, k.Select, k.Refresh, k.Install, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k RepositoriesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var RespotioriesDefaultValuesKeyHelp = RepositoriesKeyMap{
	Cancel: key.NewBinding(key.WithKeys(DefaultConfig.Repository.Cancel...), key.WithHelp(formatHelp(DefaultConfig.Repository.Cancel), "Cancel")),
}

func mergeSlices(slices ...[]string) []string {
	var out []string
	for _, s := range slices {
		out = append(out, s...)
	}
	return out
}

var RepoKeys = RepositoriesKeyMap{
	Add: key.NewBinding(key.WithKeys(DefaultConfig.Repository.Add...), key.WithHelp(formatHelp(DefaultConfig.Repository.Add), "Refresh")),
	Delete: key.NewBinding(
		key.WithKeys(DefaultConfig.Repository.Delete...),
		key.WithHelp(formatHelp(DefaultConfig.Repository.Delete), "Delete repo"),
	),
	Move: key.NewBinding(
		key.WithKeys(mergeSlices(DefaultConfig.Repository.Up, DefaultConfig.Repository.Down, DefaultConfig.Repository.Left, DefaultConfig.Repository.Right)...,
		),
		key.WithHelp(formatHelp(mergeSlices(DefaultConfig.Repository.Up, DefaultConfig.Repository.Down, DefaultConfig.Repository.Left, DefaultConfig.Repository.Right)), "Move")),
	Refresh: key.NewBinding(key.WithKeys(DefaultConfig.Repository.Refresh...), key.WithHelp(formatHelp(DefaultConfig.Repository.Refresh), "Refresh")),
	Select:  key.NewBinding(key.WithKeys(DefaultConfig.Repository.Select...), key.WithHelp(formatHelp(DefaultConfig.Repository.Select), "Select")),
	Update:  key.NewBinding(key.WithKeys(DefaultConfig.Repository.Update...), key.WithHelp(formatHelp(DefaultConfig.Repository.Update), "Update repo")),
	Install: key.NewBinding(key.WithKeys(DefaultConfig.Repository.Install...), key.WithHelp(formatHelp(DefaultConfig.Repository.Install), "Install version")),
	Up:      key.NewBinding(key.WithKeys(DefaultConfig.Repository.Up...)),
	Down:    key.NewBinding(key.WithKeys(DefaultConfig.Repository.Down...)),
	Left:    key.NewBinding(key.WithKeys(DefaultConfig.Repository.Left...)),
	Right:   key.NewBinding(key.WithKeys(DefaultConfig.Repository.Right...)),
}
