package repositories

import (
	"bytes"
	"fmt"
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
	repoNameStep int = iota
	urlStep
)

var addInputsHelper = []string{
	"Enter repo name",
	"Enter repo URL",
}

type AddModel struct {
	installStep int
	Inputs      []textinput.Model
	width       int
	height      int
	help        help.Model
	keys        keyMap
}

func InitAddModel() AddModel {
	repoName := textinput.New()
	url := textinput.New()
	inputs := []textinput.Model{repoName, url}
	m := AddModel{installStep: repoNameStep, Inputs: inputs, help: help.New(), keys: installKeys}
	return m
}

func (m AddModel) Init() tea.Cmd {
	return nil
}

func (m AddModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	cmds := make([]tea.Cmd, len(m.Inputs))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.Inputs[repoNameStep].Width = msg.Width - 6 - len(inputsHelper[0])
		m.Inputs[urlStep].Width = msg.Width - 6 - len(inputsHelper[1])
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.installStep == urlStep {
				cmds = append(cmds, m.addRepo(m.Inputs[repoNameStep].Value(), m.Inputs[urlStep].Value()))
				cmd = m.resetAllInputs()
				cmds = append(cmds, cmd)
				cmd = m.blurAllInputs()
				cmds = append(cmds, cmd)

				return m, tea.Batch(cmds...)
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
		}
	}
	return m, m.updateInputs(msg)
}

func (m AddModel) View() string {
	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	var inputs string
	for step := 0; step < len(m.Inputs); step++ {
		if step == 0 {
			inputs = fmt.Sprintf("%s %s", addInputsHelper[step], m.Inputs[step].View())
			continue
		}
		inputs = lipgloss.JoinVertical(lipgloss.Top, inputs, fmt.Sprintf("%s %s", addInputsHelper[step], m.Inputs[step].View()))
	}
	inputs = styles.ActiveStyle.Border(styles.Border).Render(inputs)
	inputs = lipgloss.JoinVertical(lipgloss.Top, inputs)
	return lipgloss.JoinVertical(lipgloss.Top, inputs, helpView)
}

func (m *AddModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m AddModel) blurAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].Blur()
	}
	return nil
}

func (m AddModel) resetAllInputs() tea.Cmd {
	for i := range m.Inputs {
		m.Inputs[i].SetValue("")
	}
	return nil
}

func (m AddModel) addRepo(repoName, url string) tea.Cmd {
	return func() tea.Msg {
		var stdout, stderr bytes.Buffer

		var cmd *exec.Cmd
		cmd = exec.Command("helm", "repo", "add", repoName, url)
		// Create the command
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// Run the command
		err := cmd.Run()
		if err != nil {
			return types.AddRepoMsg{Err: err}
		}

		return types.AddRepoMsg{Err: nil}
	}
}
