package helpers

import (
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/types"
)

func WriteAndOpenFile(content []byte, file string) tea.Cmd {
	err := os.WriteFile(file, content, 0644)

	if err != nil {

		return func() tea.Msg {

			return types.EditorFinishedMsg{Err: err}

		}

	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	c := exec.Command(editor, file)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return types.EditorFinishedMsg{Err: err}
	})
}
