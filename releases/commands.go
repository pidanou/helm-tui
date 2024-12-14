package releases

import (
	"bytes"
	"errors"
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
	cmd := exec.Command("helm", "ls", "--all-namespaces")
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.ListReleasesMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.ListReleasesMsg{Content: releases}
	}

	// remove header and empty last line
	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		updated := strings.Join(fields[3:7], " ")      // Join the parts that belong to the updated field
		remainingFields := append(fields[:3], updated) // Keep the first 3 columns and append the updated field

		// Add the rest of the fields after the updated part
		remainingFields = append(remainingFields, fields[7:]...)
		releases = append(releases, remainingFields)
	}
	return types.ListReleasesMsg{Content: releases, Err: nil}
}

func (m *Model) history() tea.Msg {
	var stdout bytes.Buffer

	if m.releaseTable.SelectedRow() == nil {
		return types.HistoryMsg{Content: nil, Err: errors.New("no release selected")}
	}

	// Create the command
	cmd := exec.Command("helm", "history", m.releaseTable.SelectedRow()[0], "--namespace", m.releaseTable.SelectedRow()[1])
	cmd.Stdout = &stdout

	// Run the command
	err := cmd.Run()
	if err != nil {
		return types.HistoryMsg{Err: err}
	}

	lines := strings.Split(stdout.String(), "\n")
	if len(lines) <= 1 {
		return types.HistoryMsg{Err: errors.New("no history found")}
	}

	// remove header and empty last line
	lines = lines[1 : len(lines)-1]

	var history []table.Row
	for _, line := range lines {
		fields := strings.Fields(line)
		updated := strings.Join(fields[1:6], " ")
		description := strings.Join(fields[9:], " ")
		remainingFields := append(fields[:1], updated)
		remainingFields = append(remainingFields, fields[6:9]...)
		remainingFields = append(remainingFields, description)
		history = append(history, remainingFields)
	}
	return types.HistoryMsg{Content: history, Err: nil}
}

func (m *Model) delete() tea.Msg {

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

func (m Model) upgrade(chart string) tea.Cmd {
	return func() tea.Msg {
		// Create the command
		cmd := exec.Command("helm", "upgrade", m.releaseTable.SelectedRow()[0], chart, "--namespace", m.releaseTable.SelectedRow()[1])
		var stout bytes.Buffer
		cmd.Stderr = &stout
		// Run the command
		err := cmd.Run()
		if err != nil {
			return types.UpgradeMsg{Err: err}
		}
		return types.UpgradeMsg{Err: nil}
	}
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
		return types.ValuesMsg{Err: errors.New("no history found")}
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
