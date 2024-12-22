package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/styles"
)

type ColumnDefinition struct {
	Title      string
	Width      int
	FlexFactor int
}

func SetTable(t *table.Model, cols []ColumnDefinition, targetWidth int) tea.Cmd {
	var columns = make([]table.Column, len(cols))
	targetWidth = targetWidth - 2 // remove the border Width
	remainingWidthAfterFixed := targetWidth
	totalFlex := 0
	for _, col := range cols {
		if col.FlexFactor != 0 {
			totalFlex += col.FlexFactor
		}
	}
	for i, col := range cols {
		if col.Width != 0 {
			columns[i] = table.Column{
				Title: col.Title,
				Width: col.Width, // remove cell padding
			}
			remainingWidthAfterFixed = remainingWidthAfterFixed - col.Width
		}
	}
	for i, col := range cols {
		if col.FlexFactor != 0 {
			columns[i] = table.Column{
				Title: col.Title,
				Width: int(remainingWidthAfterFixed*col.FlexFactor/totalFlex) - 2, // -2 to remove the cell padding
			}
		}
	}
	// fill last column with the remaning Width due to integer division
	lastCol := columns[len(columns)-1]
	totalColWidth := 0
	for _, col := range columns {
		totalColWidth += col.Width + 2 // count the cell padding
	}
	lastCol.Width = lastCol.Width + targetWidth - totalColWidth
	columns[len(columns)-1] = lastCol
	t.SetColumns(columns)
	t.SetWidth(targetWidth)
	return nil
}

func GenerateTable() table.Model {
	t := table.New()
	s := table.DefaultStyles()
	k := table.DefaultKeyMap()
	k.HalfPageUp.Unbind()
	k.PageDown.Unbind()
	k.HalfPageDown.Unbind()
	k.HalfPageDown.Unbind()
	k.GotoBottom.Unbind()
	k.GotoTop.Unbind()
	s.Header = s.Header.
		BorderStyle(styles.Border).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	t.SetStyles(s)
	t.KeyMap = k
	return t
}

func RenderTable(t table.Model, height int, width int) string {
	var topBorder string
	t.SetHeight(height)
	t.SetWidth(width)
	view := t.View()
	var baseStyle lipgloss.Style
	topBorder = styles.GenerateTopBorderWithTitle(" Releases ", t.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return lipgloss.JoinVertical(lipgloss.Left, topBorder, view)
}
