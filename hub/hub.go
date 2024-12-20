package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

type keyMap struct {
	AddRepo key.Binding
	Search  key.Binding
	Show    key.Binding
	Cancel  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.AddRepo, k.Show, k.Search, k.Cancel}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var defaultKeysHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	Show:   key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Focus table")),
}

var tableKeysHelp = keyMap{
	Show:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Show default values")),
	Search:  key.NewBinding(key.WithKeys("/"), key.WithHelp("/", "Search")),
	AddRepo: key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "Add repo")),
}

var searchKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var addRepoKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Search")),
}

var defaultValuesKeyHelp = keyMap{
	Search: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
	Cancel: key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel")),
}

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
	m := HubModel{
		searchBar:      textinput.New(),
		resultTable:    table.New(),
		defaultValueVP: viewport.New(0, 0),
		help:           help.New(),
		view:           searchView,
		repoAddInput:   textinput.New(),
	}
	m.searchBar.Placeholder = "/ to Search a package"
	m.repoAddInput.Placeholder = "Enter local repository name"
	m.resultTable.SetStyles(s)
	m.resultTable.KeyMap = k
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
			if m.resultTable.Focused() {
				if m.resultTable.SelectedRow() != nil {
					m.view = defaultValueView
					cmds = append(cmds, m.searchDefaultValue)
				}
				return m, tea.Batch(cmds...)
			}
			m.resultTable.Focus()
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

func (m HubModel) View() string {
	header := styles.InactiveStyle.Border(styles.Border).Render(m.searchBar.View())
	remainingHeight := m.height - lipgloss.Height(header) - 2 - 1 //  searchbar padding + releaseTable padding + helper
	if m.repoAddInput.Focused() {
		remainingHeight -= 3
	}
	if m.view == defaultValueView {
		m.defaultValueVP.Height = m.height - 2 - 1
		return m.renderDefaultValueView()
	}
	m.resultTable.SetHeight(remainingHeight)
	if m.searchBar.Focused() {
		header = styles.ActiveStyle.Border(styles.Border).Render(m.searchBar.View())
	}
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(defaultKeysHelp)
	if m.searchBar.Focused() {
		helpView = m.help.View(searchKeyHelp)
	}
	if m.resultTable.Focused() {
		helpView = m.help.View(tableKeysHelp)
	}
	if m.repoAddInput.Focused() {
		helpView = m.help.View(addRepoKeyHelp)
	}
	helpView += helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	style := styles.ActiveStyle.Border(styles.Border)
	if m.repoAddInput.Focused() {
		return header + "\n" + m.renderSearchTableView() + "\n" + style.Render(m.repoAddInput.View()) + "\n" + helpView
	}
	return header + "\n" + m.renderSearchTableView() + "\n" + helpView
}

func (m HubModel) renderSearchTableView() string {
	var releasesTopBorder string
	tableView := m.resultTable.View()
	var baseStyle lipgloss.Style
	releasesTopBorder = styles.GenerateTopBorderWithTitle(" Results ", m.resultTable.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	if m.resultTable.Focused() {
		releasesTopBorder = styles.GenerateTopBorderWithTitle(" Results ", m.resultTable.Width(), styles.Border, styles.ActiveStyle.Foreground(styles.HighlightColor))
		baseStyle = styles.ActiveStyle.Border(styles.Border, false, true, true)
	}
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, releasesTopBorder, tableView)
}

func (m HubModel) renderDefaultValueView() string {
	defaultValueTopBorder := styles.GenerateTopBorderWithTitle(" Default Values ", m.defaultValueVP.Width, styles.Border, styles.InactiveStyle)
	baseStyle := styles.InactiveStyle.Border(styles.Border, false, true, true)
	helperStyle := m.help.Styles.ShortSeparator
	helpView := helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return lipgloss.JoinVertical(lipgloss.Top, defaultValueTopBorder, baseStyle.Render(m.defaultValueVP.View()), m.help.View(defaultValuesKeyHelp)+helpView)
}

func (m HubModel) searchHub() tea.Msg {
	type Package struct {
		ID             string `json:"package_id"`
		NormalizedName string `json:"normalized_name"`
		Description    string `json:"description"`
		Version        string `json:"version"`
		Repository     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"repository"`
	}

	type Response struct {
		Packages []Package `json:"packages"`
	}
	url := fmt.Sprintf("https://artifacthub.io/api/v1/packages/search?offset=0&limit=20&facets=false&ts_query_web=%s&kind=0&deprecated=false&sort=relevance", m.searchBar.Value())

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.HubSearchResultMsg{Err: err}
	}

	// Set the request header
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return types.HubSearchResultMsg{Err: err}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.HubSearchResultMsg{Err: err}
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return types.HubSearchResultMsg{Err: err}
	}

	var rows []table.Row
	for _, pkg := range response.Packages {
		row := []string{}
		row = append(row, pkg.ID, pkg.Version, pkg.NormalizedName, pkg.Repository.Name, pkg.Repository.URL, pkg.Description)
		rows = append(rows, row)
	}

	return types.HubSearchResultMsg{Content: rows}
}

func (m HubModel) searchDefaultValue() tea.Msg {
	if m.resultTable.SelectedRow() == nil {
		return nil
	}

	url := fmt.Sprintf("https://artifacthub.io/api/v1/packages/%s/%s/values", m.resultTable.SelectedRow()[0], m.resultTable.SelectedRow()[1])

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.HubSearchDefaultValueMsg{Err: err}
	}

	// Set the request header
	req.Header.Set("Accept", "application/yaml")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return types.HubSearchDefaultValueMsg{Err: err}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.HubSearchDefaultValueMsg{Err: err}
	}

	return types.HubSearchDefaultValueMsg{Content: string(body)}
}

func (m HubModel) addRepo() tea.Msg {
	if m.repoAddInput.Value() == "" || m.resultTable.SelectedRow() == nil {
		return nil
	}
	cmd := exec.Command("helm", "repo", "add", m.repoAddInput.Value(), m.resultTable.SelectedRow()[4])
	err := cmd.Run()
	if err != nil {
		return types.AddRepoMsg{Err: err}
	}
	return types.AddRepoMsg{Err: nil}
}
