package styles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

var (
	InactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = tabBorderWithBottom("┘", " ", "└")
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Border(ActiveTabBorder, true)
	WindowSize        tea.WindowSizeMsg
	Border            = lipgloss.Border(lipgloss.RoundedBorder())
	HighlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	InactiveStyle     = lipgloss.NewStyle()
	ActiveStyle       = InactiveStyle.BorderForeground(HighlightColor)
)
