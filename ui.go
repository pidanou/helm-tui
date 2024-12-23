package main

import (
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/hub"
	"github.com/pidanou/helm-tui/plugins"
	"github.com/pidanou/helm-tui/releases"
	"github.com/pidanou/helm-tui/repositories"
	"github.com/pidanou/helm-tui/styles"
	"github.com/pidanou/helm-tui/types"
)

type tabIndex uint

var tabLabels = []string{"Releases", "Repositories", "Hub", "Plugins"}

const (
	releasesTab tabIndex = iota
	repositoriesTab
	hubTab
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
	m.tabContent[hubTab] = hub.InitModel()
	m.tabContent[pluginsTab] = plugins.InitModel()
	return m
}

func (m mainModel) Init() tea.Cmd {
	var cmds = []tea.Cmd{createWorkingDir, textinput.Blink}
	for _, i := range m.tabContent {
		cmds = append(cmds, i.Init())
	}
	return tea.Batch(cmds...)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case types.InitAppMsg:
		if msg.Err != nil {
			return m, tea.Quit
		}
		m.loaded = true
	case types.EditorFinishedMsg:
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
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tabContent[releasesTab], cmd = m.tabContent[releasesTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		m.tabContent[repositoriesTab], cmd = m.tabContent[repositoriesTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		m.tabContent[hubTab], cmd = m.tabContent[hubTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		m.tabContent[pluginsTab], cmd = m.tabContent[pluginsTab].Update(tea.WindowSizeMsg{Width: m.width, Height: msg.Height - lipgloss.Height(m.renderMenu())})
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "]":
			if m.state == pluginsTab {
				m.state = 0
			} else {
				m.state++
			}
		case "[":
			if m.state == releasesTab {
				m.state = pluginsTab
			} else {
				m.state--
			}
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
		case hubTab:
			m.tabContent[hubTab], cmd = m.tabContent[hubTab].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case pluginsTab:
			m.tabContent[pluginsTab], cmd = m.tabContent[pluginsTab].Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}
	}
	m.tabContent[releasesTab], cmd = m.tabContent[releasesTab].Update(msg)
	cmds = append(cmds, cmd)
	m.tabContent[repositoriesTab], cmd = m.tabContent[repositoriesTab].Update(msg)
	cmds = append(cmds, cmd)
	m.tabContent[hubTab], cmd = m.tabContent[hubTab].Update(msg)
	cmds = append(cmds, cmd)
	m.tabContent[pluginsTab], cmd = m.tabContent[pluginsTab].Update(msg)
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
			style = styles.ActiveStyle.Background(styles.HighlightColor).Padding(0, 1)
		} else {
			style = styles.InactiveStyle.Padding(0, 1)
		}
		renderedTabs = append(renderedTabs, style.Render(t))
	}
	menu := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(menu)
	return doc.String()
}

func createWorkingDir() tea.Msg {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return types.InitAppMsg{Err: err}
	}
	workingDir := path.Join(homeDir, ".helm-tui")
	err = os.MkdirAll(workingDir, 0755)
	if err != nil {
		return types.InitAppMsg{Err: err}
	}
	helpers.UserDir = workingDir
	return types.InitAppMsg{Err: nil}
}
