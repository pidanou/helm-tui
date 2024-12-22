package releases

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/components"
	"github.com/pidanou/helm-tui/types"
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
	deleting     bool
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

func InitModel() (Model, tea.Cmd) {
	table := components.GenerateTable()
	k := generateKeys()
	m := Model{releaseTable: table, historyTable: table, help: help.New(), keys: k, upgrading: false,
		installModel: InitInstallModel(), installing: false, upgradeModel: InitUpgradeModel(), deleting: false,
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
	if m.deleting {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "y":
				return m, m.delete
			case "n":
				m.deleting = false
				return m, nil
			}
		}
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
		m.deleting = false
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
		case "D":
			m.deleting = true
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
