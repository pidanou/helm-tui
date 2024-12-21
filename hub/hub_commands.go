package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/types"
)

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
