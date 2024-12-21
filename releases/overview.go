package releases

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

type selectedView int

const (
	releasesView selectedView = iota
	historyView
	notesView
	metadataView
	hooksView
	valuesView
	manifestView
)

type Model struct {
	selectedView selectedView
	keys         []keyMap
	help         help.Model
	releaseTable table.Model
	historyTable table.Model
	notesVP      viewport.Model
	metadataVP   viewport.Model
	hooksVP      viewport.Model
	valuesVP     viewport.Model
	manifestVP   viewport.Model
	installModel InstallModel
	installing   bool
	upgradeModel UpgradeModel
	upgrading    bool
	width        int
	height       int
}

var releaseCols = []components.ColumnDefinition{
	{Title: "Name", FlexFactor: 1},
	{Title: "Namespace", FlexFactor: 1},
	{Title: "Revision", Width: 10},
	{Title: "Updated", Width: 36},
	{Title: "Status", FlexFactor: 1},
	{Title: "Chart", FlexFactor: 1},
	{Title: "App version", FlexFactor: 1},
}

var historyCols = []components.ColumnDefinition{
	{Title: "Revision", FlexFactor: 1},
	{Title: "Updated", Width: 36},
	{Title: "Status", FlexFactor: 1},
	{Title: "Chart", FlexFactor: 1},
	{Title: "App version", FlexFactor: 1},
	{Title: "Description", FlexFactor: 1},
}

var menuItem = []string{
	"History",
	"Notes",
	"Metadata",
	"Hooks",
	"Values",
	"Manifest",
}

var releaseTableCache table.Model

func generateTables() (table.Model, table.Model) {
	t := table.New()
	h := table.New()
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

	t.SetStyles(s)
	h.SetStyles(s)
	t.KeyMap = k
	h.KeyMap = k
	return t, h
}

func InitModel() (tea.Model, tea.Cmd) {
	t, h := generateTables()
	k := generateKeys()
	m := Model{releaseTable: t, historyTable: h, help: help.New(), keys: k, upgrading: false,
		installModel: InitInstallModel(), installing: false, upgradeModel: InitUpgradeModel(),
	}

	m.releaseTable.Focus()
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
			}
		case types.InstallMsg:
			m.installing = false
			cmds = append(cmds, m.list)
		}
		m.installModel, cmd = m.installModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	if m.upgrading {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.upgrading = false
			}
		case types.UpgradeMsg:
			m.upgrading = false
			cmds = append(cmds, m.list)
		}
		m.upgradeModel, cmd = m.upgradeModel.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}
	switch m.selectedView {
	case releasesView:
		m.releaseTable, cmd = m.releaseTable.Update(msg)
		cmds = append(cmds, cmd)
	case historyView:
		m.historyTable, cmd = m.historyTable.Update(msg)
		cmds = append(cmds, cmd)
	case notesView:
		m.notesVP, cmd = m.notesVP.Update(msg)
		cmds = append(cmds, cmd)
	case metadataView:
		m.metadataVP, cmd = m.metadataVP.Update(msg)
		cmds = append(cmds, cmd)
	case hooksView:
		m.hooksVP, cmd = m.hooksVP.Update(msg)
		cmds = append(cmds, cmd)
	case valuesView:
		m.valuesVP, cmd = m.valuesVP.Update(msg)
		cmds = append(cmds, cmd)
	case manifestView:
		m.manifestVP, cmd = m.manifestVP.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		components.SetTable(&m.releaseTable, releaseCols, m.width)
		components.SetTable(&m.historyTable, historyCols, m.width)
		m.notesVP = viewport.New(m.width-6, 0)
		m.metadataVP = viewport.New(m.width-6, 0)
		m.hooksVP = viewport.New(m.width-6, 0)
		m.valuesVP = viewport.New(m.width-6, 0)
		m.manifestVP = viewport.New(m.width-6, 0)
		m.help.Width = msg.Width
		m.installModel, _ = m.installModel.Update(msg)
		m.upgradeModel, _ = m.upgradeModel.Update(msg)
	case types.ListReleasesMsg:
		if m.selectedView == releasesView {
			m.releaseTable.SetRows(msg.Content)
		}
		releaseTableCache = table.New(table.WithRows(msg.Content), table.WithColumns(m.releaseTable.Columns()))
		m.releaseTable, cmd = m.releaseTable.Update(msg)
		cmds = append(cmds, cmd, m.history, m.getNotes, m.getMetadata, m.getHooks, m.getValues, m.getManifest)
	case types.HistoryMsg:
		m.historyTable.SetRows(msg.Content)
		m.historyTable.SetCursor(0)
		m.historyTable, cmd = m.historyTable.Update(msg)
		cmds = append(cmds, cmd)
	case types.UpgradeMsg:
		cmds = append(cmds, m.list)
		m.selectedView = releasesView
	case types.DeleteMsg:
		cmds = append(cmds, m.list)
		m.releaseTable.SetCursor(0)
		m.selectedView = releasesView
	case types.RollbackMsg:
		cmds = append(cmds, m.history)
		m.historyTable.SetCursor(0)
	case types.NotesMsg:
		m.notesVP.SetContent(msg.Content)
		m.notesVP, cmd = m.notesVP.Update(msg)
		cmds = append(cmds, cmd)
	case types.MetadataMsg:
		m.metadataVP.SetContent(msg.Content)
		m.metadataVP, cmd = m.metadataVP.Update(msg)
		cmds = append(cmds, cmd)
	case types.HooksMsg:
		m.hooksVP.SetContent(msg.Content)
		m.hooksVP, cmd = m.hooksVP.Update(msg)
		cmds = append(cmds, cmd)
	case types.ValuesMsg:
		m.valuesVP.SetContent(msg.Content)
		m.valuesVP, cmd = m.valuesVP.Update(msg)
		cmds = append(cmds, cmd)
	case types.ManifestMsg:
		m.manifestVP.SetContent(msg.Content)
		m.manifestVP, cmd = m.manifestVP.Update(msg)
		cmds = append(cmds, cmd)
	case types.InstallMsg:
		cmds = append(cmds, m.list)

	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.installing = true
			cmd = m.installModel.Init()
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case "r":
			switch m.selectedView {
			case releasesView:
				cmds = append(cmds, m.list)
			}
		case "R":
			switch m.selectedView {
			case historyView:
				return m, m.rollback
			}
		case "d":
			return m, m.delete
		case "u":
			m.upgrading = true
			m.upgradeModel.ReleaseName = m.releaseTable.SelectedRow()[0]
			m.upgradeModel.Namespace = m.releaseTable.SelectedRow()[1]
			cmd = m.upgradeModel.Init()
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case "esc":
			m.installing = false
			m.upgrading = false
			switch m.selectedView {
			case releasesView:
			default:
				m.historyTable.SetCursor(0)
				m.selectedView = releasesView
				m.historyTable.Blur()
				m.releaseTable = releaseTableCache
			}
		case "enter", " ":
			switch m.selectedView {
			case releasesView:
				m.selectedView = historyView
				releaseTableCache = m.releaseTable
				m.releaseTable.SetHeight(3)
				m.releaseTable.SetRows([]table.Row{m.releaseTable.SelectedRow()})
				m.releaseTable.GotoTop()
				m.historyTable.Focus()
				cmds = append(cmds, m.history, m.getNotes, m.getMetadata, m.getHooks, m.getValues, m.getManifest)
			}
		case "tab":
			switch m.selectedView {
			case releasesView:
			case manifestView:
				m.selectedView = historyView
			default:
				m.selectedView++
			}
		case "shift+tab":
			switch m.selectedView {
			case releasesView:
			case historyView:
				m.selectedView = manifestView
			default:
				m.selectedView--
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var view string
	if m.installing {
		return m.installModel.View()
	}
	if m.upgrading {
		return m.upgradeModel.View()
	}

	switch m.selectedView {
	case releasesView:
		tHeight := m.height - 2 - 1 // releaseTable padding + helper
		m.releaseTable.SetHeight(tHeight)
		view = m.renderReleasesTableView()
	default:
		view = m.renderReleaseDetail()
	}

	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys[m.selectedView]) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return view + "\n" + helpView
}

func (m Model) menuView() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range menuItem {
		var style lipgloss.Style
		isFirst, isActive := i == 0, i == int(m.selectedView)-1
		if isActive {
			style = styles.ActiveTabStyle
		} else {
			style = styles.InactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row + strings.Repeat("─", m.width-lipgloss.Width(row)-1) + styles.Border.TopRight)
	return doc.String()
}

func (m Model) renderReleaseDetail() string {
	header := m.renderReleasesTableView() + "\n" + m.menuView()
	remainingHeight := m.height - lipgloss.Height(header) + lipgloss.Height(m.menuView()) - 2 - 1 // releaseTable padding + helper
	var view string
	switch m.selectedView {
	case historyView:
		m.historyTable.SetHeight(remainingHeight - 2)
		view = header + "\n" + m.renderHistoryTableView()
	case notesView:
		m.notesVP.Height = remainingHeight - 4
		view = header + "\n" + m.renderNotesView()
	case metadataView:
		m.metadataVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderMetadataView()
	case hooksView:
		m.hooksVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderHooksView()
	case valuesView:
		m.valuesVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderValuesView()
	case manifestView:
		m.manifestVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderManifestView()
	}
	return view
}

func (m Model) renderReleasesTableView() string {
	var releasesTopBorder string
	tableView := m.releaseTable.View()
	var baseStyle lipgloss.Style
	releasesTopBorder = styles.GenerateTopBorderWithTitle(" Releases ", m.releaseTable.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, releasesTopBorder, tableView)
}

func (m Model) renderHistoryTableView() string {
	tableView := m.historyTable.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Border(styles.Border).UnsetBorderTop()
	tableView = baseStyle.Render(tableView)
	return tableView
}

func (m Model) renderNotesView() string {
	view := m.notesVP.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderMetadataView() string {
	view := m.metadataVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderHooksView() string {
	view := m.hooksVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderValuesView() string {
	view := m.valuesVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderManifestView() string {
	view := m.manifestVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) list() tea.Msg {
	var stdout bytes.Buffer
	var releases = []table.Row{}

	// Create the command
	cmd := exec.Command("helm", "ls", "--all-namespaces", "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.ListReleasesMsg{Err: err}
	}
	var rls []types.Release
	err = json.Unmarshal(stdout.Bytes(), &rls)
	if err != nil {
		return types.ListReleasesMsg{Content: releases}
	}

	for _, rel := range rls {
		row := []string{rel.Name, rel.Namespace, rel.Revision, rel.Updated, rel.Status, rel.Chart, rel.AppVersion}
		releases = append(releases, row)
	}
	return types.ListReleasesMsg{Content: releases, Err: nil}
}

func (m *Model) history() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.HistoryMsg{Content: nil, Err: errors.New("no release selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "history", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1], "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.HistoryMsg{Err: err}
	}
	var history []types.History
	var rows = []table.Row{}
	err = json.Unmarshal(stdout.Bytes(), &history)
	if err != nil {
		return types.HistoryMsg{Content: rows}
	}

	for _, line := range history {
		row := []string{fmt.Sprint(line.Revision), line.Updated, line.Status, line.Chart, line.AppVersion, line.Description}
		rows = append(rows, row)
	}
	return types.HistoryMsg{Content: rows, Err: nil}
}

func (m *Model) delete() tea.Msg {

	if m.releaseTable.SelectedRow() == nil {
		return types.DeleteMsg{Err: errors.New("No release selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "uninstall", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.DeleteMsg{Err: err}
	}
	return types.DeleteMsg{Err: nil}
}

func (m Model) rollback() tea.Msg {

	// Create the command
	cmd := exec.Command("helm", "rollback", m.releaseTable.SelectedRow()[0], m.historyTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.RollbackMsg{Err: err}
	}
	return types.RollbackMsg{Err: nil}
}

func (m Model) getNotes() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.NotesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "notes", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.NotesMsg{Err: err}
	}

	return types.NotesMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getMetadata() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.NotesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "metadata", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.MetadataMsg{Err: err}
	}

	return types.MetadataMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getHooks() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.HooksMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "hooks", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.HooksMsg{Err: err}
	}

	return types.HooksMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getValues() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.ValuesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "values", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.ValuesMsg{Err: err}
	}
	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.ValuesMsg{Err: errors.New("no values found")}
	}
	lines = lines[1:]

	return types.ValuesMsg{Content: strings.Join(lines, "\n"), Err: nil}
}

func (m Model) getManifest() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.ManifestMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "manifest", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.ManifestMsg{Err: err}
	}
	return types.ManifestMsg{Content: stdout.String(), Err: nil}
}
