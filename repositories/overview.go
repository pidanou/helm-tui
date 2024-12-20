package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

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
	selectedView selectedView
	keys         []keyMap
	tables       []table.Model
	installModel InstallModel
	addModel     AddModel
	help         help.Model
	installing   bool
	adding       bool
	width        int
	height       int
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
	tableListView := table.New()
	tablePackagesView := table.New()
	tableVersionsView := table.New()

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

	tableListView.SetStyles(s)
	tableListView.KeyMap = k

	tablePackagesView.SetStyles(s)
	tablePackagesView.KeyMap = k

	tableVersionsView.SetStyles(s)
	tableVersionsView.KeyMap = k

	return tableListView, tablePackagesView, tableVersionsView
}

func InitModel() (tea.Model, tea.Cmd) {
	tables := []table.Model{}
	repoTable, tablePackagesView, tableVersionsView := generateTables()
	repoTable.Focus()
	tables = append(tables, repoTable, tablePackagesView, tableVersionsView)
	repoTable.Focus()
	keys := generateKeys()
	m := Model{
		tables:       tables,
		selectedView: listView,
		keys:         keys,
		installModel: InitInstallModel("", ""),
		addModel:     InitAddModel(),
		help:         help.New(),
		installing:   false,
		adding:       false,
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
			if msg.String() == "esc" {
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
		m.installModel.Update(msg)
		m.addModel.Update(msg)
		m.help.Width = msg.Width
	case types.ListRepoMsg:
		m.tables[listView].SetRows(msg.Content)
		m.tables[listView], cmd = m.tables[listView].Update(msg)
		cmds = append(cmds, cmd, m.searchPackages)
	case types.PackagesMsg:
		m.tables[packagesView].SetRows(msg.Content)
		m.tables[packagesView], cmd = m.tables[packagesView].Update(msg)
		cmds = append(cmds, cmd, m.searchPackageVersions)
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

	// handle key presses
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			if m.tables[packagesView].SelectedRow() != nil && m.tables[versionsView].SelectedRow() != nil {
				m.installModel.Chart = m.tables[packagesView].SelectedRow()[0]
				m.installModel.Version = m.tables[versionsView].SelectedRow()[0]
				m.installing = true
				cmd = m.installModel.Init()
				return m, cmd
			}
		case "a":
			m.adding = true
			cmd = m.addModel.Init()
			return m, cmd
		case "down", "up", "j", "k":
			switch m.selectedView {
			case listView:
				cmds = append(cmds, m.searchPackages)
			case packagesView:
				cmds = append(cmds, m.searchPackageVersions)
			}
		case "tab", "l", "enter":
			switch m.selectedView {
			case versionsView:
			default:
				m.selectedView++
			}
			m.FocusOnlyTable(m.selectedView)
		case "d":
			cmds = append(cmds, m.remove)
		case "shift+tab", "h":
			switch m.selectedView {
			case listView:
			default:
				m.selectedView--
			}
			m.FocusOnlyTable(m.selectedView)
		case "u":
			return m, m.update
		case "r":
			return m, m.list
		case "esc":
			m.installing = false
			m.adding = false
			m.selectedView = listView
		}
	}
	m.tables[listView], cmd = m.tables[listView].Update(msg)
	cmds = append(cmds, cmd)
	m.tables[packagesView], cmd = m.tables[packagesView].Update(msg)
	cmds = append(cmds, cmd)
	m.tables[versionsView], cmd = m.tables[versionsView].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	helpView := m.help.View(m.keys[m.selectedView])
	repoView := m.renderTable(m.tables[listView], " Repositories ", m.selectedView == listView)
	packagesView := m.renderTable(m.tables[packagesView], " Packages ", m.selectedView == packagesView)
	versionsView := m.renderTable(m.tables[versionsView], " Versions ", m.selectedView == versionsView)
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

func (m *Model) FocusOnlyTable(index selectedView) {
	m.tables[listView].Blur()
	m.tables[packagesView].Blur()
	m.tables[versionsView].Blur()
	m.tables[index].Focus()
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

// Commands

func (m Model) list() tea.Msg {
	var stdout bytes.Buffer
	releases := []table.Row{}

	// Create the command
	cmd := exec.Command("helm", "repo", "update")
	err := cmd.Run()

	if err != nil {
		return types.ListRepoMsg{Err: err}
	}
	cmd = exec.Command("helm", "repo", "ls")
	cmd.Stdout = &stdout

	// Run the command
	err = cmd.Run()
	if err != nil {
		return types.ListRepoMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.ListRepoMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		releases = append(releases, fields)
	}
	return types.ListRepoMsg{Content: releases, Err: nil}
}

func (m Model) update() tea.Msg {
	if m.tables[listView].SelectedRow() == nil {
		return types.UpdateRepoMsg{Err: errors.New("no repo selected")}
	}
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "update", m.tables[listView].SelectedRow()[0])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.UpdateRepoMsg{Err: err}
	}

	return types.UpdateRepoMsg{Err: nil}
}

func (m Model) remove() tea.Msg {
	if m.tables[listView].SelectedRow() == nil {
		return types.RemoveMsg{Err: errors.New("no repo selected")}
	}
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "remove", m.tables[listView].SelectedRow()[0])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.RemoveMsg{Err: err}
	}

	return types.RemoveMsg{Err: nil}
}

func (m Model) searchPackages() tea.Msg {
	var stdout bytes.Buffer
	releases := []table.Row{}
	if m.tables[listView].SelectedRow() == nil {
		return types.PackagesMsg{Content: releases, Err: errors.New("no repo selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s/", m.tables[listView].SelectedRow()[0]))
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackagesMsg{Content: releases, Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.PackagesMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		nameField := fields[0]
		releases = append(releases, table.Row{nameField})
	}
	return types.PackagesMsg{Content: releases, Err: nil}
}

func (m Model) searchPackageVersions() tea.Msg {
	var stdout bytes.Buffer
	releases := []table.Row{}
	if m.tables[packagesView].SelectedRow() == nil {
		return types.PackageVersionsMsg{Content: releases, Err: errors.New("no package selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s", m.tables[packagesView].SelectedRow()[0]), "--versions")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackageVersionsMsg{Content: releases, Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.PackageVersionsMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		allFields := strings.Fields(line)
		if len(allFields) == 0 {
			continue
		}
		joinedDescription := strings.Join(allFields[3:], " ")
		fields := allFields[1:3]
		fields = append(fields, joinedDescription)
		releases = append(releases, fields)
	}
	return types.PackageVersionsMsg{Content: releases, Err: nil}
}
