package releases

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/types"
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
	m := UpgradeModel{upgradeStep: upgradeReleaseChartStep, Inputs: inputs, help: help.New(), keys: upgradeKeys}
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
					defaultValue := true
					return m, m.openEditorWithValues(defaultValue)
				case "n":
				case "y":
					defaultValue := false
					return m, m.openEditorWithValues(defaultValue)
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
