package keymaps

// this is used for actions

var DefaultConfig = KeyConfig{
	Common: CommonKey{
		MenuNext: []string{"]"},
		MenuPrev: []string{"["},
		Exit:     []string{"ctrl+c"},
	},
	Hub: HubKey{
		StartSearch:  []string{"/"},
		SendInput:    []string{"enter"},
		StartAddRepo: []string{"a"},
		ShowValues:   []string{"?"},
		Cancel:       []string{"esc"},
	},
	Plugin: PluginKey{
		Uninstall:      []string{"U"},
		Install:        []string{"i"},
		ConfirmInstall: []string{"enter"},
		Update:         []string{"u"},
		Cancel:         []string{"esc"},
		Refresh:        []string{"r"},
	},
	Release: ReleaseKey{
		Install:  []string{"i"},
		Delete:   []string{"D"},
		Rollback: []string{"R"},
		Refresh:  []string{"r"},
		Select:   []string{"enter"},
		NextTab:  []string{"l", "right"},
		PrevTab:  []string{"h", "left"},
		Back:     []string{"esc"},
		Upgrade:  []string{"u"},
		Cancel:   []string{"esc"},
		NextStep: []string{"enter"},
	},
	Repository: RepositoryKey{
		Add:                      []string{"a"},
		Delete:                   []string{"D"},
		Refresh:                  []string{"r"},
		Up:                       []string{"k", "up"},
		Down:                     []string{"j", "down"},
		Left:                     []string{"h", "left"},
		Right:                    []string{"l", "right"},
		Update:                   []string{"u"},
		Install:                  []string{"i"},
		Select:                   []string{"enter"},
		Cancel:                   []string{"esc"},
		NextStep:                 []string{"enter"},
		ShowPackageDefaultValues: []string{"?"},
	},
}
