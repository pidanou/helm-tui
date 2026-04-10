package repositories

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/components"
	"github.com/pidanou/helm-tui/keymaps"
	"github.com/pidanou/helm-tui/types"
)

type selectedView int
type installStep int

const (
	listView selectedView = iota
	packagesView
	versionsView
)

type Model struct {
	selectedView     selectedView
	tables           []table.Model
	installModel     InstallModel
	addModel         AddModel
	help             help.Model
	installing       bool
	adding           bool
	defaultValueVP   viewport.Model
	showDefaultValue bool
	width            int
	height           int
}

var repositoryCols = []components.ColumnDefinition{
	{Title: "Name", FlexFactor: 1},
	{Title: "URL", FlexFactor: 3},
}

var packagesCols = []components.ColumnDefinition{
	{Title: "Name", FlexFactor: 1},
}

var versionsCols = []components.ColumnDefinition{
	{Title: "Chart Version", Width: 13},
	{Title: "App Version", Width: 13},
	{Title: "Description", FlexFactor: 1},
}

func InitModel() (tea.Model, tea.Cmd) {
	tables := []table.Model{}
	t := components.GenerateTable()
	repoTable := t
	tablePackagesView := t
	tableVersionsView := t
	repoTable.Focus()
	tables = append(tables, repoTable, tablePackagesView, tableVersionsView)
	repoTable.Focus()
	m := Model{
		tables:           tables,
		selectedView:     listView,
		installModel:     InitInstallModel("", ""),
		addModel:         InitAddModel(),
		help:             help.New(),
		installing:       false,
		adding:           false,
		defaultValueVP:   viewport.New(0, 0),
		showDefaultValue: false,
	}
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return m.list
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.installing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if keymaps.Contains(keymaps.DefaultConfig.Repository.Cancel, msg.String()) {
				m.installing = false
				m.installModel, cmd = m.installModel.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		case types.InstallMsg:
			m.installing = false
			cmds = append(cmds, m.list)
		}
		m.installModel, cmd = m.installModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	if m.adding {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if keymaps.Contains(keymaps.DefaultConfig.Repository.Cancel, msg.String()) {
				m.adding = false
				m.addModel, cmd = m.addModel.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		case types.AddRepoMsg:
			m.adding = false
			return m, m.list
		}
		m.addModel, cmd = m.addModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	// handle messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		components.SetTable(&m.tables[listView], repositoryCols, m.width/4)
		components.SetTable(&m.tables[packagesView], packagesCols, m.width/4)
		components.SetTable(&m.tables[versionsView], versionsCols, 2*m.width/4)
		m.defaultValueVP.Width = m.width - 2
		m.installModel.Update(msg)
		m.addModel.Update(msg)
		m.help.Width = msg.Width
	case types.ListRepoMsg:
		m.tables[listView].SetRows(msg.Content)
		m.tables[listView], cmd = m.tables[listView].Update(msg)
		cmds = append(cmds, cmd)
	case types.PackagesMsg:
		m.tables[packagesView].SetRows(msg.Content)
		m.tables[packagesView], cmd = m.tables[packagesView].Update(msg)
		cmds = append(cmds, cmd)
	case types.PackageVersionsMsg:
		m.tables[versionsView].SetRows(msg.Content)
		m.tables[versionsView], cmd = m.tables[versionsView].Update(msg)
		cmds = append(cmds, cmd)
	case types.RemoveMsg:
		cmds = append(cmds, m.list)
		m.tables[listView].SetCursor(0)
		m.selectedView = listView
	case types.InstallMsg:
		m.installing = false
	case types.AddRepoMsg:
		m.adding = false
		cmds = append(cmds, m.list)
	case types.UpdateRepoMsg:
		cmds = append(cmds, m.list)
	case types.DefaultValueMsg:
		m.defaultValueVP.SetContent(msg.Content)

	// handle key presses
	case tea.KeyMsg:
		switch {
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Install, msg.String()):
			if m.tables[packagesView].SelectedRow() != nil && m.tables[versionsView].SelectedRow() != nil {
				m.installModel.Chart = m.tables[packagesView].SelectedRow()[0]
				m.installModel.Version = m.tables[versionsView].SelectedRow()[0]
				m.installing = true
				cmd = m.installModel.Init()
				return m, cmd
			}
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Add, msg.String()):
			m.adding = true
			cmd = m.addModel.Init()
			return m, cmd
		case keymaps.Contains(keymaps.DefaultConfig.Repository.ShowPackageDefaultValues, msg.String()):
			m.showDefaultValue = true
			return m, m.getDefaultValue
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Select, msg.String()):
			switch m.selectedView {
			case listView:
				cmds = append(cmds, m.searchPackages)
			case packagesView:
				cmds = append(cmds, m.searchPackageVersions)
			}
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Right, msg.String()):
			switch m.selectedView {
			case versionsView:
			default:
				m.selectedView++
			}
			m.FocusOnlyTable(m.selectedView)
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Delete, msg.String()):
			cmds = append(cmds, m.remove)
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Left, msg.String()):
			switch m.selectedView {
			case listView:
			default:
				m.selectedView--
			}
			m.FocusOnlyTable(m.selectedView)
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Update, msg.String()):
			return m, m.update
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Refresh, msg.String()):
			return m, m.list
		case keymaps.Contains(keymaps.DefaultConfig.Repository.Cancel, msg.String()):
			m.installing = false
			m.adding = false
			m.showDefaultValue = false
			m.selectedView = listView
		}
	}
	m.tables[listView], cmd = m.tables[listView].Update(msg)
	cmds = append(cmds, cmd)
	m.tables[packagesView], cmd = m.tables[packagesView].Update(msg)
	cmds = append(cmds, cmd)
	m.tables[versionsView], cmd = m.tables[versionsView].Update(msg)
	cmds = append(cmds, cmd)
	m.defaultValueVP, cmd = m.defaultValueVP.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *Model) FocusOnlyTable(index selectedView) {
	m.tables[listView].Blur()
	m.tables[packagesView].Blur()
	m.tables[versionsView].Blur()
	m.tables[index].Focus()
}
