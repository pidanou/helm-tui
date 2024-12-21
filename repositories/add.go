package repositories

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
	addStep int
	Inputs  []textinput.Model
	width   int
	height  int
	help    help.Model
	keys    keyMap
}

func InitAddModel() AddModel {
	repoName := textinput.New()
	url := textinput.New()
	inputs := []textinput.Model{repoName, url}
	m := AddModel{addStep: repoNameStep, Inputs: inputs, help: help.New(), keys: addKeys}
	return m
}

func (m AddModel) Init() tea.Cmd {
	return m.Inputs[repoNameStep].Focus()
}

func (m AddModel) Update(msg tea.Msg) (AddModel, tea.Cmd) {
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
			if m.addStep == urlStep {
				cmds = append(cmds, m.addRepo(m.Inputs[repoNameStep].Value(), m.Inputs[urlStep].Value()))
				cmd = m.resetAllInputs()
				cmds = append(cmds, cmd)
				cmd = m.blurAllInputs()
				cmds = append(cmds, cmd)
				cmd = m.Inputs[repoNameStep].Focus()
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}

			m.addStep++

			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == int(m.addStep) {
					cmds[i] = m.Inputs[i].Focus()
					continue
				}
				m.Inputs[i].Blur()
			}

			return m, tea.Batch(cmds...)
		case "esc":
			m.addStep = 0
			for i := 0; i <= len(m.Inputs)-1; i++ {
				m.Inputs[i].Blur()
				m.Inputs[i].SetValue("")
			}
			cmd = m.Inputs[repoNameStep].Focus()
			cmds = append(cmds, cmd)
		}
	}
	cmds = append(cmds, m.updateInputs(msg))
	return m, tea.Batch(cmds...)
}
