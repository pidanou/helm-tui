package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func GenerateTopBorderWithTitle(title string, width int, border lipgloss.Border, style lipgloss.Style) string {
	var topBorder string
	// total length of top border runes, not including corners
	length := max(0, width-lipgloss.Width(title))
	leftLength := length / 2
	rightLength := max(0, length-leftLength)
	topBorder = lipgloss.JoinHorizontal(lipgloss.Left,
		border.TopLeft,
		strings.Repeat(border.Top, leftLength),
		title,
		strings.Repeat(border.Top, rightLength),
		border.TopRight,
	)
	return style.Render(topBorder)
}
