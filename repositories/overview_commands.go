package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/pidanou/helmtui/types"
)

func (m Model) list() tea.Msg {
	var stdout bytes.Buffer
	repositories := []table.Row{}

	// Create the command
	cmd := exec.Command("helm", "repo", "update")
	err := cmd.Run()

	if err != nil {
		return types.ListRepoMsg{Err: err}
	}
	cmd = exec.Command("helm", "repo", "ls", "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err = cmd.Run()
	if err != nil {
		return types.ListRepoMsg{Err: err}
	}

	var repos []types.Repository
	err = json.Unmarshal(stdout.Bytes(), &repos)
	if err != nil {
		return []string{}
	}

	for _, repo := range repos {
		row := []string{repo.Name, repo.URL}
		repositories = append(repositories, row)
	}
	return types.ListRepoMsg{Content: repositories, Err: nil}
}

func (m Model) update() tea.Msg {
	if m.tables[listView].SelectedRow() == nil {
		return types.UpdateRepoMsg{Err: errors.New("no repo selected")}
	}
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "update", m.tables[listView].SelectedRow()[0])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.UpdateRepoMsg{Err: err}
	}

	return types.UpdateRepoMsg{Err: nil}
}

func (m Model) remove() tea.Msg {
	if m.tables[listView].SelectedRow() == nil {
		return types.RemoveMsg{Err: errors.New("no repo selected")}
	}
	var stdout bytes.Buffer

	// Create the command
	cmd := exec.Command("helm", "repo", "remove", m.tables[listView].SelectedRow()[0])
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
	if m.tables[listView].SelectedRow() == nil {
		return types.PackagesMsg{Content: releases, Err: errors.New("no repo selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s/", m.tables[listView].SelectedRow()[0]), "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackagesMsg{Content: releases, Err: err}
	}
	var pkgs []types.Pkg
	err = json.Unmarshal(stdout.Bytes(), &pkgs)
	if err != nil {
		return []string{}
	}

	for _, pkg := range pkgs {
		releases = append(releases, table.Row{pkg.Name})
	}
	return types.PackagesMsg{Content: releases, Err: nil}
}

func (m Model) searchPackageVersions() tea.Msg {
	var stdout bytes.Buffer
	versions := []table.Row{}
	if m.tables[packagesView].SelectedRow() == nil {
		return types.PackageVersionsMsg{Content: versions, Err: errors.New("no package selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "search", "repo", fmt.Sprintf("%s", m.tables[packagesView].SelectedRow()[0]), "--versions", "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.PackageVersionsMsg{Content: versions, Err: err}
	}
	var pkgs []types.Pkg
	err = json.Unmarshal(stdout.Bytes(), &pkgs)
	if err != nil {
		return []string{}
	}

	for _, pkg := range pkgs {
		versions = append(versions, table.Row{pkg.Version, pkg.AppVersion, pkg.Description})
	}

	return types.PackageVersionsMsg{Content: versions, Err: nil}
}

func (m Model) getDefaultValue() tea.Msg {
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "show", "values", fmt.Sprintf("%s", m.tables[packagesView].SelectedRow()[0]), "--version", m.tables[versionsView].SelectedRow()[0])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.DefaultValueMsg{Content: "Unable to get default values"}
	}
	return types.DefaultValueMsg{Content: stdout.String()}
}
