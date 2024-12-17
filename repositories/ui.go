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
type installStep int

const (
	listView selectedView = iota
	packagesView
	versionsView
)

const (
	nameStep installStep = iota
	namespaceStep
	valuesStep
	confirmStep
)

type Model struct {
	selectedView      selectedView
	keys              []keyMap
	repositoriesTable table.Model
	packagesTable     table.Model
	versionsTable     table.Model
	inputs            []textinput.Model
	help              help.Model
	installing        bool
	installStep       installStep
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
	name := textinput.New()
	namespace := textinput.New()
	value := textinput.New()
	confirm := textinput.New()
	inputs := []textinput.Model{name, namespace, value, confirm}
	m := Model{repositoriesTable: repoTable,
		packagesTable: packagesTable,
		versionsTable: versionsTable,
		selectedView:  listView,
		keys:          keys,
		help:          help.New(),
		inputs:        inputs,
		installing:    false,
	}
	m.inputs[nameStep].Placeholder = "Enter release name"
	m.inputs[namespaceStep].Placeholder = "Enter namespace (empty for current)"
	m.inputs[valuesStep].Placeholder = "Edit default values ? y/n"
	m.inputs[confirmStep].Placeholder = "Enter to install"
	return m, nil
}

func (m Model) Init() tea.Cmd {
	return m.list
}

func (m Model) handleInstallSteps(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case types.EditorFinishedMsg:
		m.installStep++
		for i := 0; i <= len(m.inputs)-1; i++ {
			if i == int(m.installStep) {
				cmds[i] = m.inputs[i].Focus()
				continue
			}
			m.inputs[i].Blur()
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.installStep == confirmStep {
				cmd = m.blurAllInputs()
				cmds = append(cmds, cmd)

				cmd = m.installPackage()
				cmds = append(cmds, cmd)
				m.installing = false

				return m, tea.Batch(cmds...)
			}

			if m.installStep == valuesStep {
				switch m.inputs[valuesStep].Value() {
				case "y":
					return m, m.openEditorDefaultValues()
				case "n":
				}
			}

			m.installStep++

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == int(m.installStep) {
					cmds[i] = m.inputs[i].Focus()
					continue
				}
				m.inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		case "esc":
			for i := 0; i <= len(m.inputs)-1; i++ {
				m.inputs[i].Blur()
				m.inputs[i].SetValue("")
			}
			m.installing = false
		}
	}
	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	if m.installing {
		return m.handleInstallSteps(msg)
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
		m.help.Width = msg.Width
		m.inputs[nameStep].Width = msg.Width - 6
		m.inputs[namespaceStep].Width = msg.Width - 6
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
		m.cleanValueFile()
		cmds = append(cmds, m.cleanValueFile)

	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.installing = true
			m.installStep = nameStep
			cmd = m.inputs[nameStep].Focus()
			cmds = append(cmds, cmd)
			// return m, tea.Batch(cmds...)
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
		case "left", "h":
			switch m.selectedView {
			case listView:
			default:
				m.selectedView--
			}
		case "R":
			return m, m.remove
		case "esc":
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
	var inputs string
	for step := 0; step < len(m.inputs); step++ {
		if step == 0 {
			inputs = m.inputs[step].View()
			continue
		}
		inputs = lipgloss.JoinVertical(lipgloss.Top, inputs, m.inputs[step].View())
	}
	inputs = styles.ActiveStyle.Border(styles.Border).Render(inputs)
	if m.installing {
		view = lipgloss.JoinVertical(lipgloss.Top, inputs)
		return lipgloss.JoinVertical(lipgloss.Top, view, helpView)
	}
	return view + "\n" + strings.Repeat(" ", m.width-lipgloss.Width(helpView)) + helpView
}

func (m Model) renderTable(table table.Model, title string, active bool) string {
	var topBorder string
	table.SetHeight(m.height - 3)
	if m.installing {
		table.SetHeight(m.height - 3 - lipgloss.Height(m.inputs[nameStep].View()) - lipgloss.Height(m.inputs[namespaceStep].View()) - 2)
	}
	tableView := table.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.InactiveStyle)
	if active && !m.installing {
		topBorder = styles.GenerateTopBorderWithTitle(title, table.Width(), styles.Border, styles.ActiveStyle.Foreground(styles.HighlightColor))
		baseStyle = styles.ActiveStyle.Border(styles.Border, false, true, true)
	}
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, topBorder, tableView)
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) blurAllInputs() tea.Cmd {
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	return nil
}
