package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tea "github.com/charmbracelet/bubbletea"
)

// TestInitAddModel verifies that the AddModel initializes correctly.
func TestInitAddModel(t *testing.T) {
	model := InitAddModel()

	assert.Equal(t, repoNameStep, model.installStep, "Initial installStep should be repoNameStep")
	assert.Equal(t, 2, len(model.Inputs), "AddModel should have 2 inputs")
	assert.True(t, model.Inputs[0].Focused(), "First input should be focused on initialization")
}

// TestAddModelEnterKey verifies that pressing Enter advances the install step.
func TestAddModelEnterKey(t *testing.T) {
	model := InitAddModel()
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)

	assert.Equal(t, urlStep, updatedModel.installStep, "installStep should advance to urlStep after pressing Enter")
	assert.True(t, updatedModel.Inputs[urlStep].Focused(), "Second input should be focused after pressing Enter")
}

// TestAddModelEscKey verifies that pressing Esc resets the install step and clears inputs.
func TestAddModelEscKey(t *testing.T) {
	model := InitAddModel()
	model.Inputs[repoNameStep].SetValue("test-repo")
	model.Inputs[urlStep].SetValue("http://example.com")

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)

	assert.Equal(t, repoNameStep, updatedModel.installStep, "installStep should be reset to repoNameStep")
	for _, input := range updatedModel.Inputs {
		assert.Empty(t, input.Value(), "All inputs should be cleared after pressing Esc")
	}
}