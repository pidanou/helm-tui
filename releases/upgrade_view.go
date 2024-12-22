package releases

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/styles"
)

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
