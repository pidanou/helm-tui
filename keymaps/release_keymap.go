package keymaps

import "github.com/charmbracelet/bubbles/key"

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type ReleaseKeyMap struct {
	Install  key.Binding
	Delete   key.Binding
	Rollback key.Binding
	Refresh  key.Binding
	Select   key.Binding
	MoveTab  key.Binding
	Back     key.Binding
	Upgrade  key.Binding
	NextStep key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k ReleaseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.NextStep, k.Install, k.Delete, k.Upgrade, k.Select, k.Refresh, k.Rollback, k.MoveTab, k.Back}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k ReleaseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var ReleasesKeys = ReleaseKeyMap{
	Refresh: key.NewBinding(key.WithKeys(DefaultConfig.Release.Refresh...), key.WithHelp(formatHelp(DefaultConfig.Release.Refresh), "Refresh")),
	Select:  key.NewBinding(key.WithKeys(DefaultConfig.Release.Select...), key.WithHelp(formatHelp(DefaultConfig.Release.Select), "Details")),
	Install: key.NewBinding(key.WithKeys(DefaultConfig.Release.Install...), key.WithHelp(formatHelp(DefaultConfig.Release.Install), "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys(DefaultConfig.Release.Delete...),
		key.WithHelp(formatHelp(DefaultConfig.Release.Delete), "Delete release"),
	),
	Upgrade: key.NewBinding(key.WithKeys(DefaultConfig.Release.Upgrade...), key.WithHelp(formatHelp(DefaultConfig.Release.Upgrade), "Upgrade release")),
}

var HistoryKeys = ReleaseKeyMap{
	Install:  key.NewBinding(key.WithKeys(DefaultConfig.Release.Install...), key.WithHelp(formatHelp(DefaultConfig.Release.Install), "Install new release")),
	Rollback: key.NewBinding(key.WithKeys(DefaultConfig.Release.Rollback...), key.WithHelp(formatHelp(DefaultConfig.Release.Rollback), "Rollback to revision")),
	Delete: key.NewBinding(
		key.WithKeys(DefaultConfig.Release.Delete...),
		key.WithHelp(formatHelp(DefaultConfig.Release.Delete), "Delete release"),
	),
	Upgrade: key.NewBinding(key.WithKeys(DefaultConfig.Release.Upgrade...), key.WithHelp(formatHelp(DefaultConfig.Release.Upgrade), "Upgrade release")),
	MoveTab: key.NewBinding(key.WithKeys(append(DefaultConfig.Release.NextTab, DefaultConfig.Release.PrevTab...)...), key.WithHelp(formatHelp(append(DefaultConfig.Release.NextTab, DefaultConfig.Release.PrevTab...)), "Navigate tabs")),
	Back:    key.NewBinding(key.WithKeys(DefaultConfig.Release.Back...), key.WithHelp(formatHelp(DefaultConfig.Release.Back), "Back")),
}

var ReadOnlyKeys = ReleaseKeyMap{
	Install: key.NewBinding(key.WithKeys(DefaultConfig.Release.Install...), key.WithHelp(formatHelp(DefaultConfig.Release.Install), "Install new release")),
	Delete: key.NewBinding(
		key.WithKeys(DefaultConfig.Release.Delete...),
		key.WithHelp(formatHelp(DefaultConfig.Release.Delete), "Delete release"),
	),
	Upgrade: key.NewBinding(key.WithKeys(DefaultConfig.Release.Upgrade...), key.WithHelp(formatHelp(DefaultConfig.Release.Upgrade), "Upgrade release")),
	MoveTab: key.NewBinding(key.WithKeys(append(DefaultConfig.Release.NextTab, DefaultConfig.Release.PrevTab...)...), key.WithHelp(formatHelp(append(DefaultConfig.Release.NextTab, DefaultConfig.Release.PrevTab...)), "Navigate tabs")),
	Back:    key.NewBinding(key.WithKeys(DefaultConfig.Release.Back...), key.WithHelp(formatHelp(DefaultConfig.Release.Back), "Back")),
}

func GenerateReleaseKeys() []ReleaseKeyMap {
	return []ReleaseKeyMap{ReleasesKeys, HistoryKeys, ReadOnlyKeys, ReadOnlyKeys, ReadOnlyKeys, ReadOnlyKeys, ReadOnlyKeys}

}

// Only used for installation form
var ReleaseInstallKeys = ReleaseKeyMap{
	Back:     key.NewBinding(key.WithKeys(DefaultConfig.Release.Back...), key.WithHelp(formatHelp(DefaultConfig.Release.Back), "Cancel")),
	NextStep: key.NewBinding(key.WithKeys(DefaultConfig.Release.NextStep...), key.WithHelp(formatHelp(DefaultConfig.Release.NextStep), "Next step")),
}
