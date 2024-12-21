package releases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

const (
	upgradeReleaseChartStep int = iota
	upgradeReleaseVersionStep
	upgradeReleaseValuesStep
	upgradeReleaseConfirmStep
)

var upgradeInputsHelper = []string{
	"Enter a chart name or chart directory (absolute path)",
	"Version (empty for latest)",
	"Edit values yes/no/use default ? y/n/d",
	"Confirm ? enter/esc",
}

type UpgradeModel struct {
	ReleaseName string
	Namespace   string
	upgradeStep int
	Chart       string
	Version     string
	Inputs      []textinput.Model
	width       int
	height      int
	help        help.Model
	keys        keyMap
	tag         int
}

func InitUpgradeModel() UpgradeModel {
	chart := textinput.New()
	version := textinput.New()
	value := textinput.New()
	confirm := textinput.New()
	inputs := []textinput.Model{chart, version, value, confirm}
	m := UpgradeModel{upgradeStep: upgradeReleaseChartStep, Inputs: inputs, help: help.New(), keys: installKeys}
	m.Inputs[upgradeReleaseChartStep].ShowSuggestions = true
	m.Inputs[upgradeReleaseVersionStep].ShowSuggestions = true
	return m
}

func (m UpgradeModel) Init() tea.Cmd {
	return m.Inputs[0].Focus()
}

func (m UpgradeModel) Update(msg tea.Msg) (UpgradeModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, len(m.Inputs))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Inputs[upgradeReleaseChartStep].Width = msg.Width - 6 - len(upgradeInputsHelper[0])
		m.Inputs[upgradeReleaseValuesStep].Width = msg.Width - 6 - len(upgradeInputsHelper[1])
	case types.UpgradeMsg:
		m.upgradeStep = 0
		if m.Namespace == "" {
			m.Namespace = "default"
		}
		folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, m.Namespace, m.ReleaseName)
		cmds = append(cmds, m.cleanValueFile(folder), m.blurAllInputs(), m.resetAllInputs())
		return m, tea.Batch(cmds...)

	case types.DebounceEndMsg:
		if msg.Tag == m.tag {
			if m.Inputs[upgradeReleaseChartStep].Focused() {
				m.Inputs[upgradeReleaseChartStep].SetSuggestions(m.searchLocalPackage())
			}
			if m.Inputs[upgradeReleaseVersionStep].Focused() {
				m.Inputs[upgradeReleaseVersionStep].SetSuggestions(m.searchLocalPackageVersion())
			}
		}
	case types.EditorFinishedMsg:
		m.upgradeStep++
		for i := 0; i <= len(m.Inputs)-1; i++ {
			if i == int(m.upgradeStep) {
				cmds[i] = m.Inputs[i].Focus()
				continue
			}
			m.Inputs[i].Blur()
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		m.tag++
		switch msg.String() {
		case "enter":
			if m.upgradeStep == upgradeReleaseConfirmStep {
				m.upgradeStep = 0
				cmd = m.blurAllInputs()
				cmds = append(cmds, cmd, m.upgrade)

				return m, tea.Batch(cmds...)
			}

			if m.upgradeStep == upgradeReleaseValuesStep {
				switch m.Inputs[upgradeReleaseValuesStep].Value() {
				case "d":
					return m, m.openEditorDefaultValues()
				case "n":
				case "y":
					return m, m.openEditorLastValues()
				default:
					return m, nil
				}
			}

			m.upgradeStep++

			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == int(m.upgradeStep) {
					cmds[i] = m.Inputs[i].Focus()
					continue
				}
				m.Inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		case "esc":
			m.upgradeStep = 0
			for i := 0; i <= len(m.Inputs)-1; i++ {
				m.Inputs[i].Blur()
				m.Inputs[i].SetValue("")
			}
		default:
			return m, tea.Batch(m.updateInputs(msg), tea.Tick(debounce, func(_ time.Time) tea.Msg {
				return types.DebounceEndMsg{Tag: m.tag}
			}))

		}
	}
	cmd = m.updateInputs(msg)
	return m, cmd
}

func (m UpgradeModel) View() string {
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	if m.Inputs[upgradeReleaseChartStep].Focused() {
		helpView = m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.SuggestionInputKeyMap) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	}
	var Inputs string
	for step := 0; step < len(m.Inputs); step++ {
		if step == 0 {
			Inputs = fmt.Sprintf("%s %s", upgradeInputsHelper[step], m.Inputs[step].View())
			continue
		}
		Inputs = lipgloss.JoinVertical(lipgloss.Top, Inputs, fmt.Sprintf("%s %s", upgradeInputsHelper[step], m.Inputs[step].View()))
	}
	Inputs = styles.ActiveStyle.Border(styles.Border).Render(Inputs)
	Inputs = lipgloss.JoinVertical(lipgloss.Top, Inputs)
	return lipgloss.JoinVertical(lipgloss.Top, Inputs, helpView)
}

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
