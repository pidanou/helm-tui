package releases

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/helpers"
	"github.com/pidanou/helmtui/styles"
)

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
