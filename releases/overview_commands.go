package releases

import (
	"bytes"
	"encoding/json"
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
	var releases = []table.Row{}

	// Create the command
	cmd := exec.Command("helm", "ls", "--all-namespaces", "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.ListReleasesMsg{Err: err}
	}
	var rls []types.Release
	err = json.Unmarshal(stdout.Bytes(), &rls)
	if err != nil {
		return types.ListReleasesMsg{Content: releases}
	}

	for _, rel := range rls {
		row := []string{rel.Name, rel.Namespace, rel.Revision, rel.Updated, rel.Status, rel.Chart, rel.AppVersion}
		releases = append(releases, row)
	}
	return types.ListReleasesMsg{Content: releases, Err: nil}
}

func (m *Model) history() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.HistoryMsg{Content: nil, Err: errors.New("no release selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "history", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1], "--output", "json")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.HistoryMsg{Err: err}
	}
	var history []types.History
	var rows = []table.Row{}
	err = json.Unmarshal(stdout.Bytes(), &history)
	if err != nil {
		return types.HistoryMsg{Content: rows}
	}

	for _, line := range history {
		row := []string{fmt.Sprint(line.Revision), line.Updated, line.Status, line.Chart, line.AppVersion, line.Description}
		rows = append(rows, row)
	}
	return types.HistoryMsg{Content: rows, Err: nil}
}

func (m *Model) delete() tea.Msg {

	if m.releaseTable.SelectedRow() == nil {
		return types.DeleteMsg{Err: errors.New("No release selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "uninstall", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.DeleteMsg{Err: err}
	}
	return types.DeleteMsg{Err: nil}
}

func (m Model) rollback() tea.Msg {

	// Create the command
	cmd := exec.Command("helm", "rollback", m.releaseTable.SelectedRow()[0], m.historyTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.RollbackMsg{Err: err}
	}
	return types.RollbackMsg{Err: nil}
}

func (m Model) getNotes() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.NotesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "notes", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.NotesMsg{Err: err}
	}

	return types.NotesMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getMetadata() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.NotesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "metadata", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.MetadataMsg{Err: err}
	}

	return types.MetadataMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getHooks() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.HooksMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "hooks", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.HooksMsg{Err: err}
	}

	return types.HooksMsg{Content: stdout.String(), Err: nil}
}

func (m Model) getValues() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.ValuesMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "values", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.ValuesMsg{Err: err}
	}
	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.ValuesMsg{Err: errors.New("no values found")}
	}
	lines = lines[1:]

	return types.ValuesMsg{Content: strings.Join(lines, "\n"), Err: nil}
}

func (m Model) getManifest() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.ManifestMsg{Err: errors.New("no release selected")}
	}

	cmd := exec.Command("helm", "get", "manifest", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return types.ManifestMsg{Err: err}
	}
	return types.ManifestMsg{Content: stdout.String(), Err: nil}
}
