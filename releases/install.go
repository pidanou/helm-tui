package releases

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	installChartReleaseNameStep int = iota
	installChartNameStep
	installChartVersionStep
	installChartNamespaceStep
	installChartValuesStep
	installChartConfirmStep
)

var installInputsHelper = []string{
	"Enter release name",
	"Enter chart",
	"Enter chart version (empty for latest)",
	"Enter namespace (empty for default)",
	"Edit default values ? y/n",
	"Enter to install",
}

const debounce = 500 * time.Millisecond

type InstallModel struct {
	installStep int
	Chart       string
	Version     string
	Inputs      []textinput.Model
	width       int
	height      int
	help        help.Model
	keys        keyMap
	tag         int
}

func InitInstallModel() InstallModel {
	chart := textinput.New()
	version := textinput.New()
	name := textinput.New()
	namespace := textinput.New()
	value := textinput.New()
	confirm := textinput.New()
	inputs := []textinput.Model{name, chart, version, namespace, value, confirm}
	m := InstallModel{installStep: installChartReleaseNameStep, Inputs: inputs, help: help.New(), keys: installKeys}
	m.Inputs[installChartNameStep].ShowSuggestions = true
	m.Inputs[installChartVersionStep].ShowSuggestions = true
	return m
}

func (m InstallModel) Init() tea.Cmd {
	return m.Inputs[0].Focus()
}

func (m InstallModel) Update(msg tea.Msg) (InstallModel, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, len(m.Inputs))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.Inputs[installChartReleaseNameStep].Width = msg.Width - 6 - len(installInputsHelper[0])
	case types.EditorFinishedMsg:
		m.installStep++
		for i := 0; i <= len(m.Inputs)-1; i++ {
			if i == int(m.installStep) {
				cmds[i] = m.Inputs[i].Focus()
				continue
			}
			m.Inputs[i].Blur()
		}
		return m, tea.Batch(cmds...)
	case types.InstallMsg:
		m.installStep = 0
		releaseName := m.Inputs[installChartReleaseNameStep].Value()
		namespace := m.Inputs[installChartNamespaceStep].Value()
		if namespace == "" {
			namespace = "default"
		}
		folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, namespace, releaseName)
		cmds = append(cmds, m.cleanValueFile(folder), m.blurAllInputs(), m.resetAllInputs())

		return m, tea.Batch(cmds...)
	case types.DebounceEndMsg:
		if msg.Tag == m.tag {
			if m.Inputs[installChartNameStep].Focused() {
				m.Inputs[installChartNameStep].SetSuggestions(m.searchLocalPackage())
			}
			if m.Inputs[installChartVersionStep].Focused() {
				m.Inputs[installChartVersionStep].SetSuggestions(m.searchLocalPackageVersion())
			}
		}
	case tea.KeyMsg:
		m.tag++
		switch msg.String() {
		case "enter":
			if m.installStep == installChartConfirmStep {
				m.installStep = 0

				cmd = m.installPackage(m.Inputs[installChartValuesStep].Value())
				cmds = append(cmds, cmd)

				return m, tea.Batch(cmds...)
			}

			if m.installStep == installChartValuesStep {
				switch m.Inputs[installChartValuesStep].Value() {
				case "y":
					return m, m.openEditorDefaultValues()
				case "n":
				default:
					return m, nil
				}
			}

			m.installStep++

			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == int(m.installStep) {
					cmds[i] = m.Inputs[i].Focus()
					continue
				}
				m.Inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		case "esc":
			m.installStep = 0
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
	return m, m.updateInputs(msg)
}

func (m InstallModel) View() string {
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	if m.Inputs[installChartNameStep].Focused() {
		helpView = m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.SuggestionInputKeyMap) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	}
	var inputs string
	for step := 0; step < len(m.Inputs); step++ {
		if step == 0 {
			inputs = fmt.Sprintf("%s %s", installInputsHelper[step], m.Inputs[step].View())
			continue
		}
		inputs = lipgloss.JoinVertical(lipgloss.Top, inputs, fmt.Sprintf("%s %s", installInputsHelper[step], m.Inputs[step].View()))
	}
	inputs = styles.ActiveStyle.Border(styles.Border).Render(inputs)
	inputs = lipgloss.JoinVertical(lipgloss.Top, inputs)
	return lipgloss.JoinVertical(lipgloss.Top, inputs, helpView)
}

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

func (m InstallModel) searchLocalPackage() []string {
	if m.Inputs[installChartNameStep].Value() == "" {
		return []string{}
	}
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", m.Inputs[installChartNameStep].Value())
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	results := strings.Split(stdout.String(), "\n")
	results = results[1 : len(results)-1]

	var suggestions []string
	for _, row := range results {
		chart := strings.Fields(row)
		suggestions = append(suggestions, chart[0])
	}

	m.Inputs[installChartNameStep].SetSuggestions(suggestions)
	return suggestions
}

func (m InstallModel) searchLocalPackageVersion() []string {
	var stdout bytes.Buffer
	cmd := exec.Command("helm", "search", "repo", "--regexp", "\v"+m.Inputs[installChartNameStep].Value()+"\v", "--versions")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return []string{}
	}
	results := strings.Split(stdout.String(), "\n")
	results = results[1 : len(results)-1]

	var suggestions []string
	for _, row := range results {
		chart := strings.Fields(row)
		suggestions = append(suggestions, chart[1])
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
