package repositories

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/styles"
)

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
