package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/releases"
	"github.com/pidanou/helmtui/repositories"
	"github.com/pidanou/helmtui/styles"
)

type tabIndex uint

var tabLabels = []string{"[1] releases", "[2] chart", "[3] repositories", "[4] plugins"}

const (
	releasesTab tabIndex = iota
	chartTab
	repositoriesTab
	pluginsTab
)

type mainModel struct {
	state      tabIndex
	index      int
	width      int
	height     int
	tabs       []string
	tabContent []tea.Model
	loaded     bool
}

func newModel(tabs []string) mainModel {
	m := mainModel{state: releasesTab, tabs: tabs, tabContent: make([]tea.Model, len(tabs)), loaded: false}
	m.tabContent[releasesTab], _ = releases.InitModel()
	m.tabContent[repositoriesTab], _ = repositories.InitModel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.tabContent[releasesTab].Init())
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tabContent[releasesTab], cmd = m.tabContent[releasesTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		m.tabContent[repositoriesTab], cmd = m.tabContent[repositoriesTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		m.loaded = true
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "1":
			m.state = releasesTab
		// case "2":
		// 	m.state = chart
		case "3":
			m.state = repositoriesTab
			cmds = append(cmds, m.tabContent[repositoriesTab].Init())
			// case "4":
			// 	m.state = plugins
		}
		switch m.state {
		case releasesTab:
			m.tabContent[releasesTab], cmd = m.tabContent[releasesTab].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case repositoriesTab:
			m.tabContent[repositoriesTab], cmd = m.tabContent[repositoriesTab].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}
	m.tabContent[releasesTab], cmd = m.tabContent[releasesTab].Update(msg)
	cmds = append(cmds, cmd)
	m.tabContent[repositoriesTab], cmd = m.tabContent[repositoriesTab].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	doc := strings.Builder{}
	if !m.loaded || len(m.tabContent) == 0 {
		return "loading..."
	}
	doc.WriteString(m.renderMenu())
	doc.WriteString("\n")
	doc.WriteString(m.tabContent[m.state].View())
	return doc.String()
}

func (m mainModel) renderMenu() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range m.tabs {
		var style lipgloss.Style
		isActive := i == int(m.state)
		if isActive {
			style = styles.ActiveStyle.BorderStyle(styles.Border)
		} else {
			style = styles.InactiveStyle.BorderStyle(styles.Border)
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	menu := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(menu)
	return doc.String()
}
