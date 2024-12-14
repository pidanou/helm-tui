package repositories

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

type selectedView int

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
	releaseNameInput  textinput.Model
	namespaceInput    textinput.Model
	help              help.Model
	installing        bool
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
	ti := textinput.New()
	ns := textinput.New()
	m := Model{repositoriesTable: repoTable,
		packagesTable:    packagesTable,
		versionsTable:    versionsTable,
		selectedView:     listView,
		keys:             keys,
		help:             help.New(),
		releaseNameInput: ti,
		namespaceInput:   ns,
		installing:       false,
	}
	m.repositoriesTable.Focus()
	m.releaseNameInput.Placeholder = "Enter release name"
	m.namespaceInput.Placeholder = "Enter namespace (empty for default)"
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
			switch msg.String() {
			case "enter":
				if m.releaseNameInput.Focused() {
					m.releaseNameInput.Blur()
					cmd = m.namespaceInput.Focus()
					return m, cmd
				} else if m.namespaceInput.Focused() {
					cmd = m.installPackage(m.releaseNameInput.Value(), m.namespaceInput.Value())
					cmds = append(cmds, cmd)
					m.releaseNameInput.SetValue("")
					m.releaseNameInput.Blur()
					m.namespaceInput.SetValue("")
					m.namespaceInput.Blur()
					m.installing = false
				}
				return m, tea.Batch(cmds...)
			case "esc":
				m.releaseNameInput.SetValue("")
				m.namespaceInput.SetValue("")
				m.releaseNameInput.Blur()
				m.namespaceInput.Blur()
				m.installing = false
			}
		}
		m.releaseNameInput, cmd = m.releaseNameInput.Update(msg)
		cmds = append(cmds, cmd)
		m.namespaceInput, cmd = m.namespaceInput.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		components.SetTable(&m.repositoriesTable, repositoryCols, m.width/4)
		components.SetTable(&m.packagesTable, packagesCols, m.width/4)
		components.SetTable(&m.versionsTable, versionsCols, 2*m.width/4)
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
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.installing = true
			cmd = m.releaseNameInput.Focus()
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case "down", "up":
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
		case "left", "h":
			switch m.selectedView {
			case listView:
			default:
				m.selectedView--
			}
		case "R":
			return m, m.remove
		case "esc":
			m.packagesTable.SetCursor(-1)
			m.selectedView = listView
		}
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
		view = lipgloss.JoinVertical(lipgloss.Top, view, m.releaseNameInput.View(), m.namespaceInput.View())
	}
	return view + "\n" + strings.Repeat(" ", m.width-lipgloss.Width(helpView)) + helpView
}

func (m Model) renderTable(table table.Model, title string, active bool) string {
	var topBorder string
	table.SetHeight(m.height - 3)
	if m.installing {
		table.SetHeight(m.height - 3 - lipgloss.Height(m.releaseNameInput.View()) - lipgloss.Height(m.namespaceInput.View()))
	}
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
