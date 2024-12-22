package repositories

import (
	"bytes"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/types"
)

func (m *AddModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m AddModel) blurAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
	return nil
}

func (m AddModel) resetAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}
	return nil
}

func (m AddModel) addRepo(repoName, url string) tea.Cmd {
	return func() tea.Msg {
		var stdout, stderr bytes.Buffer

		var cmd *exec.Cmd
		cmd = exec.Command("helm", "repo", "add", repoName, url)
		// Create the command
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Run the command
		err := cmd.Run()
		if err != nil {
			return types.AddRepoMsg{Err: err}
		}

		return types.AddRepoMsg{Err: nil}
	}
}
