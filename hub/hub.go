package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/components"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

type HubModel struct {
	searchBar   textinput.Model
	resultTable table.Model
	help        help.Model
	width       int
	height      int
}

var resultsCols = []components.ColumnDefinition{
	{Title: "id", Width: 0},
	{Title: "Package", FlexFactor: 1},
	{Title: "Repository", FlexFactor: 1},
	{Title: "URL", FlexFactor: 3},
	{Title: "Description", FlexFactor: 3},
}

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
		searchBar:   textinput.New(),
		resultTable: table.New(),
		help:        help.New(),
	}
	m.searchBar.Placeholder = "Search a package"
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
	case types.HubSearchResultMsg:
		m.resultTable.SetRows(msg.Content)
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			m.resultTable.Blur()
			cmds = append(cmds, m.searchBar.Focus())
			return m, tea.Batch(cmds...)
		case "enter":
			if m.searchBar.Focused() {
				m.searchBar.Blur()
				m.resultTable.Focus()
				cmds = append(cmds, m.searchHub)
			}
		}
	}
	m.searchBar, cmd = m.searchBar.Update(msg)
	cmds = append(cmds, cmd)
	m.resultTable, cmd = m.resultTable.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m HubModel) View() string {
	header := styles.InactiveStyle.Border(styles.Border).Render(m.searchBar.View())
	remainingHeight := m.height - lipgloss.Height(header) - 2 - 1 //  searchbar padding + releaseTable padding + helper
	m.resultTable.SetHeight(remainingHeight)
	if m.searchBar.Focused() {
		header = styles.ActiveStyle.Border(styles.Border).Render(m.searchBar.View())
	}
	helperStyle := m.help.Styles.ShortSeparator
	helpView := helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return header + "\n" + m.renderSearchTableView() + "\n" + helpView
}

func (m HubModel) renderSearchTableView() string {
	var releasesTopBorder string
	tableView := m.resultTable.View()
	var baseStyle lipgloss.Style
	releasesTopBorder = styles.GenerateTopBorderWithTitle(" Results ", m.resultTable.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, releasesTopBorder, tableView)
}

func (m HubModel) searchHub() tea.Msg {
	type Package struct {
		ID             string `json:"package_id"`
		NormalizedName string `json:"normalized_name"`
		Description    string `json:"description"`
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
		helpers.Println("Error creating request:", err)
		return types.HubSearchResultMsg{Err: err}
	}

	// Set the request header
	req.Header.Set("Accept", "application/json")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		helpers.Println("Error making request:", err)
		return types.HubSearchResultMsg{Err: err}
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.Println("Error reading response:", err)
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
		row = append(row, pkg.ID, pkg.NormalizedName, pkg.Repository.Name, pkg.Repository.URL, pkg.Description)
		rows = append(rows, row)
	}

	return types.HubSearchResultMsg{Content: rows}
}
