package hub

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/types"
)

type HubModel struct {
	searchBar      textinput.Model
	resultTable    table.Model
	defaultValueVP viewport.Model
	repoAddInput   textinput.Model
	help           help.Model
	width          int
	height         int
	view           int
}

var resultsCols = []components.ColumnDefinition{
	{Title: "id", Width: 0},
	{Title: "version", Width: 0},
	{Title: "Package", FlexFactor: 1},
	{Title: "Repository", FlexFactor: 1},
	{Title: "URL", FlexFactor: 3},
	{Title: "Description", FlexFactor: 3},
}

const (
	searchView int = iota
	defaultValueView
)

func InitModel() tea.Model {
	resultTable := components.GenerateTable()
	m := HubModel{
		searchBar:      textinput.New(),
		resultTable:    resultTable,
		defaultValueVP: viewport.New(0, 0),
		help:           help.New(),
		view:           searchView,
		repoAddInput:   textinput.New(),
	}
	m.searchBar.Placeholder = "/ to Search a package"
	m.repoAddInput.Placeholder = "Enter local repository name"
	return m
}

func (m HubModel) Init() tea.Cmd {
	return nil
}

func (m HubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.searchBar.Width = msg.Width - 5 // -2 for border, -1 for input chevron
		components.SetTable(&m.resultTable, resultsCols, m.width)
		m.defaultValueVP.Width = m.width - 2
		m.repoAddInput.Width = m.width - 5
	case types.HubSearchResultMsg:
		m.resultTable.SetRows(msg.Content)
	case types.HubSearchDefaultValueMsg:
		m.defaultValueVP.SetContent(msg.Content)
	case types.AddRepoMsg:
		m.repoAddInput.SetValue("")
		m.repoAddInput.Blur()
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if !m.repoAddInput.Focused() && !m.searchBar.Focused() {
				m.resultTable.Blur()
				m.searchBar.Blur()
				cmds = append(cmds, m.repoAddInput.Focus())
				return m, tea.Batch(cmds...)
			}
		case "/":
			if m.view == searchView {
				m.resultTable.Blur()
				cmds = append(cmds, m.searchBar.Focus())
				return m, tea.Batch(cmds...)
			}
		case "enter":
			if m.repoAddInput.Focused() {
				cmds = append(cmds, m.addRepo)
				return m, tea.Batch(cmds...)
			}
			if m.searchBar.Focused() {
				m.searchBar.Blur()
				m.resultTable.Focus()
				cmds = append(cmds, m.searchHub)
				return m, tea.Batch(cmds...)
			}
			m.resultTable.Focus()
		case "v":
			if m.resultTable.Focused() {
				if m.resultTable.SelectedRow() != nil {
					m.view = defaultValueView
					cmds = append(cmds, m.searchDefaultValue)
				}
				return m, tea.Batch(cmds...)
			}
		case "esc":
			m.view = searchView
			m.repoAddInput.Blur()
			m.defaultValueVP.GotoTop()
		}
	}
	m.searchBar, cmd = m.searchBar.Update(msg)
	cmds = append(cmds, cmd)
	m.resultTable, cmd = m.resultTable.Update(msg)
	cmds = append(cmds, cmd)
	m.defaultValueVP, cmd = m.defaultValueVP.Update(msg)
	cmds = append(cmds, cmd)
	m.repoAddInput, cmd = m.repoAddInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
