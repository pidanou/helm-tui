package repositories

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/types"
)

func (m Model) list() tea.Msg {
	var stdout bytes.Buffer
	var releases []table.Row

	// Create the command
	cmd := exec.Command("helm", "repo", "update")
	err := cmd.Run()

	if err != nil {
		return types.ListRepoMsg{Err: err}
	}
	cmd = exec.Command("helm", "repo", "ls")
	cmd.Stdout = &stdout

	// Run the command
	err = cmd.Run()
	if err != nil {
		return types.ListRepoMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.ListRepoMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		releases = append(releases, fields)
	}
	return types.ListRepoMsg{Content: releases, Err: nil}
}

func (m Model) remove() tea.Msg {
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "remove", m.repositoriesTable.SelectedRow()[0])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.RemoveMsg{Err: err}
	}

	return types.RemoveMsg{Err: nil}
}

func (m Model) searchPackages() tea.Msg {
	var stdout bytes.Buffer
	var releases []table.Row

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s/", m.repositoriesTable.SelectedRow()[0]))
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackagesMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.PackagesMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		nameField := fields[0]
		releases = append(releases, table.Row{nameField})
	}
	return types.PackagesMsg{Content: releases, Err: nil}
}

func (m Model) searchPackageVersions() tea.Msg {
	var stdout bytes.Buffer
	var releases []table.Row

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s", m.packagesTable.SelectedRow()[0]), "--versions")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackageVersionsMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.PackageVersionsMsg{Content: releases}
	}

	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		allFields := strings.Fields(line)
		if len(allFields) == 0 {
			continue
		}
		joinedDescription := strings.Join(allFields[3:], " ")
		fields := allFields[1:3]
		fields = append(fields, joinedDescription)
		releases = append(releases, fields)
	}
	return types.PackageVersionsMsg{Content: releases, Err: nil}
}

func (m Model) installPackage() tea.Cmd {
	releaseName := m.inputs[nameStep].Value()
	namespace := m.inputs[namespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s", namespace, releaseName)
	file := fmt.Sprintf("./%s/values.yaml", folder)
	return func() tea.Msg {
		var stdout, stderr bytes.Buffer

		var cmd *exec.Cmd
		// Create the command
		if m.inputs[valuesStep].Value() == "y" {
			cmd = exec.Command("helm", "install", releaseName, m.packagesTable.SelectedRow()[0], "--version", m.versionsTable.SelectedRow()[0], "--values", file, "--namespace", namespace, "--create-namespace")
		} else {
			cmd = exec.Command("helm", "install", releaseName, m.packagesTable.SelectedRow()[0], "--version", m.versionsTable.SelectedRow()[0], "--namespace", namespace, "--create-namespace")
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

func (m Model) openEditorDefaultValues() tea.Cmd {
	var stdout, stderr bytes.Buffer
	releaseName := m.inputs[nameStep].Value()
	namespace := m.inputs[namespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s", namespace, releaseName)
	file := fmt.Sprintf("./%s/values.yaml", folder)
	packageName := m.packagesTable.SelectedRow()[0]
	version := m.versionsTable.SelectedRow()[0]

	_ = os.MkdirAll(folder, 0755)
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

func (m Model) cleanValueFile() tea.Msg {
	releaseName := m.inputs[nameStep].Value()
	namespace := m.inputs[namespaceStep].Value()
	if namespace == "" {
		namespace = "default"
	}
	folder := fmt.Sprintf("%s/%s", namespace, releaseName)
	_ = os.RemoveAll(folder)
	return nil
}
