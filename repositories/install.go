package repositories

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
	"github.com/pidanou/helmtui/types"
)

const (
	nameStep installStep = iota
	namespaceStep
	valuesStep
	confirmStep
)

var inputsHelper = []string{
	"Enter release name",
	"Enter namespace (empty for default)",
	"Edit default values ? y/n",
	"Enter to install",
}

type InstallModel struct {
	installStep installStep
	Chart       string
	Version     string
	Inputs      []textinput.Model
	width       int
	height      int
	help        help.Model
	keys        keyMap
}

func InitInstallModel(chart, version string) InstallModel {
	name := textinput.New()
	namespace := textinput.New()
	value := textinput.New()
	confirm := textinput.New()
	inputs := []textinput.Model{name, namespace, value, confirm}
	m := InstallModel{installStep: nameStep, Inputs: inputs, help: help.New(), Chart: chart, Version: version, keys: installKeys}
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
		m.Inputs[nameStep].Width = msg.Width - 6 - len(inputsHelper[0])
		m.Inputs[namespaceStep].Width = msg.Width - 6 - len(inputsHelper[1])
		m.Inputs[valuesStep].Width = msg.Width - 6 - len(inputsHelper[2])
		m.Inputs[confirmStep].Width = msg.Width - 6 - len(inputsHelper[3])
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
		releaseName := m.Inputs[nameStep].Value()
		namespace := m.Inputs[namespaceStep].Value()
		if namespace == "" {
			namespace = "default"
		}
		folder := fmt.Sprintf("%s/%s/%s", helpers.UserDir, namespace, releaseName)
		cmds = append(cmds, m.cleanValueFile(folder), m.blurAllInputs(), m.resetAllInputs(), m.Inputs[nameStep].Focus())

		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.installStep == confirmStep {
				m.installStep = 0

				cmd = m.installPackage(m.Inputs[valuesStep].Value())
				cmds = append(cmds, cmd)

				m.Inputs[confirmStep].Blur()
				cmd = m.Inputs[nameStep].Focus()
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}

			if m.installStep == valuesStep {
				switch m.Inputs[valuesStep].Value() {
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
			cmds = append(cmds, m.Inputs[repoNameStep].Focus())
		}
	}
	cmds = append(cmds, m.updateInputs(msg))
	return m, tea.Batch(cmds...)
}

func (m InstallModel) View() string {
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	var inputs string
	for step := 0; step < len(m.Inputs); step++ {
		if step == 0 {
			inputs = fmt.Sprintf("%s %s", inputsHelper[step], m.Inputs[step].View())
			continue
		}
		inputs = lipgloss.JoinVertical(lipgloss.Top, inputs, fmt.Sprintf("%s %s", inputsHelper[step], m.Inputs[step].View()))
	}
	inputs = styles.ActiveStyle.Border(styles.Border).Render(inputs)
	inputs = lipgloss.JoinVertical(lipgloss.Top, inputs)
	return lipgloss.JoinVertical(lipgloss.Top, "\n", "Installing "+m.Chart+" "+m.Version, "\n", inputs, helpView)
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
