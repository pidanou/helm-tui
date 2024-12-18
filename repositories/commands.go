package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helmtui/types"
)

func (m Model) list() tea.Msg {
	var stdout bytes.Buffer
	releases := []table.Row{}

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

func (m Model) update() tea.Msg {
	if m.repositoriesTable.SelectedRow() == nil {
		return types.UpdateRepoMsg{Err: errors.New("no repo selected")}
	}
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "update", m.repositoriesTable.SelectedRow()[0])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.UpdateRepoMsg{Err: err}
	}

	return types.UpdateRepoMsg{Err: nil}
}

func (m Model) remove() tea.Msg {
	if m.repositoriesTable.SelectedRow() == nil {
		return types.RemoveMsg{Err: errors.New("no repo selected")}
	}
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
	releases := []table.Row{}
	if m.repositoriesTable.SelectedRow() == nil {
		return types.PackagesMsg{Content: releases, Err: errors.New("no repo selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s/", m.repositoriesTable.SelectedRow()[0]))
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackagesMsg{Content: releases, Err: err}
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
	releases := []table.Row{}
	if m.packagesTable.SelectedRow() == nil {
		return types.PackageVersionsMsg{Content: releases, Err: errors.New("no package selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s", m.packagesTable.SelectedRow()[0]), "--versions")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackageVersionsMsg{Content: releases, Err: err}
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
