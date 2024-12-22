package releases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/types"
)

func (m InstallModel) installPackage(mode string) tea.Cmd {
	chartName := m.Inputs[installChartNameStep].Value()
	version := m.Inputs[installChartVersionStep].Value()
	releaseName := m.Inputs[installChartReleaseNameStep].Value()
	namespace := m.Inputs[installChartNamespaceStep].Value()
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
			cmd = exec.Command("helm", "install", releaseName, chartName, "--version", version, "--values", file, "--namespace", namespace, "--create-namespace")
		} else {
			cmd = exec.Command("helm", "install", releaseName, chartName, "--version", version, "--namespace", namespace, "--create-namespace")
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
	releaseName := m.Inputs[installChartReleaseNameStep].Value()
	namespace := m.Inputs[installChartNamespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, namespace, releaseName)
	_ = os.MkdirAll(folder, 0755)
	file := fmt.Sprintf("%s/values.yaml", folder)
	packageName := m.Inputs[installChartNameStep].Value()
	version := m.Inputs[installChartVersionStep].Value()

	cmd := exec.Command("helm", "show", "values", packageName, "--version", version)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return func() tea.Msg { return types.EditorFinishedMsg{Err: err} }
	}
	return helpers.WriteAndOpenFile(stdout.Bytes(), file)
}

func (m InstallModel) searchLocalPackage() []string {
	if m.Inputs[installChartNameStep].Value() == "" {
		return []string{}
	}
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", m.Inputs[installChartNameStep].Value(), "--output", "json")
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

func (m InstallModel) searchLocalPackageVersion() []string {
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", "--regexp", "\v"+m.Inputs[installChartNameStep].Value()+"\v", "--versions", "--output", "json")
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

	m.Inputs[installChartNameStep].SetSuggestions(suggestions)
	return suggestions
}

func (m InstallModel) cleanValueFile(folder string) tea.Cmd {
	return func() tea.Msg {
		_ = os.RemoveAll(folder)
		return nil
	}
}
