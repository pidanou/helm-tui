package releases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/types"
)

func (m *UpgradeModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text Inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m UpgradeModel) blurAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
	return nil
}

func (m UpgradeModel) resetAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}
	return nil
}

func (m UpgradeModel) upgrade() tea.Msg {
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, m.Namespace, m.ReleaseName)
	_ = os.MkdirAll(folder, 0755)
	file := fmt.Sprintf("%s/values.yaml", folder)
	var cmd *exec.Cmd
	if m.Inputs[upgradeReleaseValuesStep].Value() == "y" || m.Inputs[upgradeReleaseValuesStep].Value() == "d" {
		cmd = exec.Command("helm", "upgrade", m.ReleaseName, m.Inputs[upgradeReleaseChartStep].Value(), "--values", file, "--namespace", m.Namespace)
	} else {
		cmd = exec.Command("helm", "upgrade", m.ReleaseName, m.Inputs[upgradeReleaseChartStep].Value(), "--namespace", m.Namespace)
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.UpgradeMsg{Err: err}
	}
	return types.UpgradeMsg{Err: nil}
}

func (m UpgradeModel) openEditorDefaultValues() tea.Cmd {
	var stdout, stderr bytes.Buffer
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, m.Namespace, m.ReleaseName)
	_ = os.MkdirAll(folder, 0755)
	file := fmt.Sprintf("%s/values.yaml", folder)
	packageName := m.Inputs[upgradeReleaseChartStep].Value()
	version := m.Inputs[upgradeReleaseVersionStep].Value()

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

func (m UpgradeModel) openEditorLastValues() tea.Cmd {
	var stdout, stderr bytes.Buffer
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, m.Namespace, m.ReleaseName)
	_ = os.MkdirAll(folder, 0755)
	file := fmt.Sprintf("%s/values.yaml", folder)

	cmd := exec.Command("helm", "get", "values", m.ReleaseName, "--namespace", m.Namespace)
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

func (m UpgradeModel) searchLocalPackage() []string {
	if m.Inputs[upgradeReleaseChartStep].Value() == "" {
		return []string{}
	}
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", m.Inputs[upgradeReleaseChartStep].Value(), "--output", "json")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	var pkgs []types.Pkg
	err = json.Unmarshal(stdout.Bytes(), &pkgs)
	if err != nil {
		return []string{}
	}
	var suggestions []string
	for _, p := range pkgs {
		suggestions = append(suggestions, p.Name)
	}

	return suggestions
}

func (m UpgradeModel) searchLocalPackageVersion() []string {
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", "--regexp", "\v"+m.Inputs[upgradeReleaseChartStep].Value()+"\v", "--versions", "--output", "json")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	var pkgs []types.Pkg
	err = json.Unmarshal(stdout.Bytes(), &pkgs)
	if err != nil {
		return []string{}
	}

	var suggestions []string
	for _, pkg := range pkgs {
		suggestions = append(suggestions, pkg.Version)
	}

	return suggestions
}

func (m UpgradeModel) cleanValueFile(folder string) tea.Cmd {
	return func() tea.Msg {
		_ = os.RemoveAll(folder)
		return nil
	}
}
