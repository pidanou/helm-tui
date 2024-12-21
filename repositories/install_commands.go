package repositories

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/types"
)

func (m *InstallModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m InstallModel) blurAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
	return nil
}

func (m InstallModel) resetAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}
	return nil
}

func (m InstallModel) installPackage(mode string) tea.Cmd {
	releaseName := m.Inputs[nameStep].Value()
	namespace := m.Inputs[namespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, namespace, releaseName)
	file := fmt.Sprintf("%s/values.yaml", folder)
	return func() tea.Msg {
		var stdout, stderr bytes.Buffer

		var cmd *exec.Cmd
		// Create the command
		if mode == "y" {
			cmd = exec.Command("helm", "install", releaseName, m.Chart, "--version", m.Version, "--values", file, "--namespace", namespace, "--create-namespace")
		} else {
			cmd = exec.Command("helm", "install", releaseName, m.Chart, "--version", m.Version, "--namespace", namespace, "--create-namespace")
		}
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Run the command
		err := cmd.Run()
		if err != nil {
			return types.InstallMsg{Err: err}
		}

		return types.InstallMsg{Err: nil}
	}
}

func (m InstallModel) openEditorDefaultValues() tea.Cmd {
	var stdout, stderr bytes.Buffer
	releaseName := m.Inputs[nameStep].Value()
	namespace := m.Inputs[namespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, namespace, releaseName)
	_ = os.MkdirAll(folder, 0755)
	file := fmt.Sprintf("%s/values.yaml", folder)
	packageName := m.Chart
	version := m.Version

	cmd := exec.Command("helm", "show", "values", packageName, "--version", version)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return func() tea.Msg { return types.EditorFinishedMsg{Err: err} }
	}
	err = os.WriteFile(file, stdout.Bytes(), 0644)
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

func (m InstallModel) cleanValueFile(folder string) tea.Cmd {
	return func() tea.Msg {
		_ = os.RemoveAll(folder)
		return nil
	}
}
