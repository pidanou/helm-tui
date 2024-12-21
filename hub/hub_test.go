package hub

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/types"
)

// Test the initialization of the HubModel
func TestInitModel(t *testing.T) {
	model := InitModel()
	if _, ok := model.(HubModel); !ok {
		t.Error("InitModel did not return a HubModel")
	}
}

// Test the Update function with a WindowSizeMsg
func TestHubModelUpdateWindowSizeMsg(t *testing.T) {
	model := InitModel().(HubModel)
	msg := tea.WindowSizeMsg{Width: 100, Height: 40}
	updatedModel, _ := model.Update(msg)

	hubModel, ok := updatedModel.(HubModel)
	if !ok {
		t.Error("Expected HubModel after update")
	}

	if hubModel.width != 100 || hubModel.height != 40 {
		t.Errorf("Expected width=100 and height=40, got width=%d and height=%d", hubModel.width, hubModel.height)
	}
}

// Test the searchHub function with an empty search input
func TestSearchHub(t *testing.T) {
	model := HubModel{}
	model.searchBar.SetValue("")

	msg := model.searchHub()
	if msg == nil {
		t.Error("Expected a non-nil message from searchHub")
	}

	if _, ok := msg.(types.HubSearchResultMsg); !ok {
		t.Error("Expected message of type HubSearchResultMsg")
	}
}

// Test the addRepo function with empty input
func TestAddRepo(t *testing.T) {
	model := HubModel{
		repoAddInput: textinput.New(),
		resultTable: table.New(
			table.WithColumns([]table.Column{
				{Title: "ID", Width: 10},
				{Title: "Version", Width: 10},
				{Title: "Package", Width: 20},
				{Title: "Repository", Width: 20},
				{Title: "URL", Width: 30},
				{Title: "Description", Width: 40},
			}),
			table.WithRows([]table.Row{
				{"1", "1.0.0", "example-package", "example-repo", "https://example.com", "An example package"},
			}),
		),
	}

	// Select the first row in the table
	model.resultTable.SetCursor(0)

	// Set a value for repoAddInput
	model.repoAddInput.SetValue("example-repo")

	msg := model.addRepo()
	if msg == nil {
		t.Error("Expected a non-nil message when adding a repo with valid input")
	}

	if _, ok := msg.(types.AddRepoMsg); !ok {
		t.Error("Expected message of type AddRepoMsg")
	}
}
