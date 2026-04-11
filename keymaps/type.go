package keymaps

type KeyConfig struct {
	Common     CommonKey     `json:"common"`
	Hub        HubKey        `json:"hub"`
	Plugin     PluginKey     `json:"plugin"`
	Release    ReleaseKey    `json:"release"`
	Repository RepositoryKey `json:"repository"`
}

type CommonKey struct {
	MenuNext []string `json:"menuNext"`
	MenuPrev []string `json:"menuPrev"`
	Exit     []string `json:"exit"`
}

type HubKey struct {
	StartSearch  []string `json:"startSearch"`
	StartAddRepo []string `json:"startAddRepo"`
	SendInput    []string `json:"sendInput"`
	ShowValues   []string `json:"showValues"`
	Cancel       []string `json:"cancel"`
}

type PluginKey struct {
	Uninstall      []string `json:"uninstall"`
	Install        []string `json:"install"`
	ConfirmInstall []string `json:"confirmInstall"`
	Update         []string `json:"update"`
	Cancel         []string `json:"cancel"`
	Refresh        []string `json:"refresh"`
}

type ReleaseKey struct {
	Install  []string `json:"install"`
	Delete   []string `json:"delete"`
	Rollback []string `json:"rollback"`
	Refresh  []string `json:"refresh"`
	Select   []string `json:"select"`
	NextTab  []string `json:"nextTab"`
	PrevTab  []string `json:"prevTab"`
	Back     []string `json:"back"`
	Upgrade  []string `json:"upgrade"`
	Cancel   []string `json:"cancel"`
	NextStep []string `json:"nextStep"`
}

type RepositoryKey struct {
	Add                      []string `json:"add"`
	Delete                   []string `json:"delete"`
	Refresh                  []string `json:"refresh"`
	Up                       []string `json:"up"`
	Down                     []string `json:"down"`
	Left                     []string `json:"left"`
	Right                    []string `json:"right"`
	Update                   []string `json:"update"`
	Install                  []string `json:"install"`
	Select                   []string `json:"select"`
	Cancel                   []string `json:"cancel"`
	NextStep                 []string `json:"nextStep"`
	ShowPackageDefaultValues []string `json:"showPackageDefaultValues"`
}
