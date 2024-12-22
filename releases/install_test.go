package releases

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/types"
	"github.com/stretchr/testify/assert"
)

// TestInitInstallModel verifies that the InstallModel initializes correctly.
func TestInitInstallModel(t *testing.T) {
	model := InitInstallModel()

	assert.Equal(t, installChartReleaseNameStep, model.installStep, "Initial installStep should be installChartReleaseNameStep")
	assert.Equal(t, 6, len(model.Inputs), "InstallModel should have 6 inputs")
}

// TestInstallModelEnterKey verifies that the Enter key advances the install step.
func TestInstallModelEnterKey(t *testing.T) {
	model := InitInstallModel()
	msg := tea.KeyMsg{Type: tea.KeyEnter}

	updatedModel, _ := model.Update(msg)

	assert.Equal(t, installChartNameStep, updatedModel.installStep, "installStep should advance to installChartNameStep after pressing Enter")
	assert.True(t, updatedModel.Inputs[installChartNameStep].Focused(), "Next input should be focused after pressing Enter")
}

// TestInstallModelEscKey verifies that the Esc key resets the install step and clears inputs.
func TestInstallModelEscKey(t *testing.T) {
	model := InitInstallModel()
	model.Inputs[installChartReleaseNameStep].SetValue("test-release")

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := model.Update(msg)

	assert.Equal(t, installChartReleaseNameStep, updatedModel.installStep, "installStep should be reset to installChartReleaseNameStep")
	for _, input := range updatedModel.Inputs {
		assert.Empty(t, input.Value(), "All inputs should be cleared after pressing Esc")
	}
}

// TestInstallMsgHandling verifies that the model resets after handling an InstallMsg.
func TestInstallMsgHandling(t *testing.T) {
	model := InitInstallModel()
	model.Inputs[installChartReleaseNameStep].SetValue("test-release")

	msg := types.InstallMsg{}
	updatedModel, _ := model.Update(msg)

	assert.Equal(t, installChartReleaseNameStep, updatedModel.installStep, "installStep should be reset to installChartReleaseNameStep after InstallMsg")
	for _, input := range updatedModel.Inputs {
		assert.Empty(t, input.Value(), "All inputs should be cleared after InstallMsg")
	}
}
