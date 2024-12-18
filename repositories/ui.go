package repositories

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

type selectedView int
type installStep int

const (
	listView selectedView = iota
	packagesView
	versionsView
)

type Model struct {
	selectedView      selectedView
	keys              []keyMap
	repositoriesTable table.Model
	packagesTable     table.Model
	versionsTable     table.Model
	installModel      InstallModel
	addModel          AddModel
	help              help.Model
	installing        bool
	adding            bool
	width             int
	height            int
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

func generateTables() (table.Model, table.Model, table.Model) {
	repositoriesTable := table.New()
	packagesTable := table.New()
	versionsTable := table.New()

	s := table.DefaultStyles()
	k := table.DefaultKeyMap()
	k.HalfPageUp.Unbind()
	k.PageDown.Unbind()
	k.HalfPageDown.Unbind()
	k.HalfPageDown.Unbind()
	k.GotoBottom.Unbind()
	k.GotoTop.Unbind()
	s.Header = s.Header.
		BorderStyle(styles.Border).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	repositoriesTable.SetStyles(s)
	repositoriesTable.KeyMap = k

	packagesTable.SetStyles(s)
	packagesTable.KeyMap = k

	versionsTable.SetStyles(s)
	versionsTable.KeyMap = k

	return repositoriesTable, packagesTable, versionsTable
}

func InitModel() (tea.Model, tea.Cmd) {
	repoTable, packagesTable, versionsTable := generateTables()
	repoTable.Focus()
	packagesTable.Focus()
	versionsTable.Focus()
	keys := generateKeys()
	m := Model{repositoriesTable: repoTable,
		packagesTable: packagesTable,
		versionsTable: versionsTable,
		selectedView:  listView,
		keys:          keys,
		installModel:  InitInstallModel("", ""),
		addModel:      InitAddModel(),
		help:          help.New(),
		installing:    false,
		adding:        false,
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
			if msg.String() == "esc" {
				m.installing = false
				m.installModel.Update(msg)
				return m, nil
			}
		case types.InstallMsg:
			m.installing = false
			mod, cmd := m.installModel.Update(msg)
			cmds = append(cmds, cmd)
			m.installModel = mod.(InstallModel)
			return m, cmd
		}
		mod, cmd := m.installModel.Update(msg)
		cmds = append(cmds, cmd)
		m.installModel = mod.(InstallModel)
		return m, tea.Batch(cmds...)
	}
	if m.adding {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.adding = false
				m.addModel.Update(msg)
				return m, nil
			}
		case types.AddRepoMsg:
			m.adding = false
			return m, m.list
		}
		mod, cmd := m.addModel.Update(msg)
		cmds = append(cmds, cmd)
		m.addModel = mod.(AddModel)
		return m, tea.Batch(cmds...)
	}

	// handle table updates
	switch m.selectedView {
	case listView:
		m.repositoriesTable, cmd = m.repositoriesTable.Update(msg)
		cmds = append(cmds, cmd)
	case packagesView:
		m.packagesTable, cmd = m.packagesTable.Update(msg)
		cmds = append(cmds, cmd)
	case versionsView:
		m.versionsTable, cmd = m.versionsTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	// handle messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		components.SetTable(&m.repositoriesTable, repositoryCols, m.width/4)
		components.SetTable(&m.packagesTable, packagesCols, m.width/4)
		components.SetTable(&m.versionsTable, versionsCols, 2*m.width/4)
		m.installModel.Update(msg)
		m.addModel.Update(msg)
		m.help.Width = msg.Width
	case types.ListRepoMsg:
		m.repositoriesTable.SetRows(msg.Content)
		m.repositoriesTable, cmd = m.repositoriesTable.Update(msg)
		cmds = append(cmds, cmd, m.searchPackages)
	case types.PackagesMsg:
		m.packagesTable.SetRows(msg.Content)
		m.packagesTable, cmd = m.packagesTable.Update(msg)
		cmds = append(cmds, cmd, m.searchPackageVersions)
	case types.PackageVersionsMsg:
		m.versionsTable.SetRows(msg.Content)
		m.versionsTable, cmd = m.versionsTable.Update(msg)
		cmds = append(cmds, cmd)
	case types.RemoveMsg:
		cmds = append(cmds, m.list)
		m.repositoriesTable.SetCursor(0)
		m.selectedView = listView
	case types.InstallMsg:
		m.installing = false
	case types.AddRepoMsg:
		m.adding = false
		cmds = append(cmds, m.list)
	case types.UpdateRepoMsg:
		cmds = append(cmds, m.list)

	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			if m.packagesTable.SelectedRow() != nil && m.versionsTable.SelectedRow() != nil {
				m.installModel.Chart = m.packagesTable.SelectedRow()[0]
				m.installModel.Version = m.versionsTable.SelectedRow()[0]
				cmd = m.installModel.Inputs[0].Focus()
				m.installing = true
				cmds = append(cmds, cmd)
			}
		case "a":
			cmd = m.addModel.Inputs[0].Focus()
			m.adding = true
			cmds = append(cmds, cmd)
		case "down", "up", "j", "k":
			switch m.selectedView {
			case listView:
				cmds = append(cmds, m.searchPackages)
			case packagesView:
				cmds = append(cmds, m.searchPackageVersions)
			}
		case "right", "l", "enter":
			switch m.selectedView {
			case versionsView:
			default:
				m.selectedView++
			}
		case "d":
			cmds = append(cmds, m.remove)
		case "left", "h":
			switch m.selectedView {
			case listView:
			default:
				m.selectedView--
			}
		case "u":
			return m, m.update
		case "r":
			return m, m.list
		case "esc":
			m.installing = false
			m.adding = false
			m.selectedView = listView
		}
	default:
		mod, cmd := m.installModel.Update(msg)
		cmds = append(cmds, cmd)
		m.installModel = mod.(InstallModel)
		mod, cmd = m.addModel.Update(msg)
		cmds = append(cmds, cmd)
		m.addModel = mod.(AddModel)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpView := m.help.View(m.keys[m.selectedView])
	repoView := m.renderTable(m.repositoriesTable, " Repositories ", m.selectedView == listView)
	packagesView := m.renderTable(m.packagesTable, " Packages ", m.selectedView == packagesView)
	versionsView := m.renderTable(m.versionsTable, " Versions ", m.selectedView == versionsView)
	view := lipgloss.JoinHorizontal(lipgloss.Top, repoView, packagesView, versionsView)
	if m.installing {
		return m.installModel.View()
	}
	if m.adding {
		return m.addModel.View()
	}
	helperStyle := m.help.Styles.ShortSeparator
	return view + "\n" + helpView + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
}

func (m Model) renderTable(table table.Model, title string, active bool) string {
	var topBorder string
	table.SetHeight(m.height - 3)
	tableView := table.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.InactiveStyle)
	if active {
		topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.ActiveStyle.Foreground(styles.HighlightColor))
		baseStyle = styles.ActiveStyle.Border(styles.Border, false, true, true)
	}
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, topBorder, tableView)
}
