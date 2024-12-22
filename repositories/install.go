package repositories

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/types"
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
		m.Inputs[nameStep].Width = msg.Width - 5 - len(inputsHelper[0])
		m.Inputs[namespaceStep].Width = msg.Width - 5 - len(inputsHelper[1])
		m.Inputs[valuesStep].Width = msg.Width - 5 - len(inputsHelper[2])
		m.Inputs[confirmStep].Width = msg.Width - 5 - len(inputsHelper[3])
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
