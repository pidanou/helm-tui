package releases

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helmtui/styles"
)

type ColumnDefinition struct {
	title      string
	width      int
	flexFactor int
}

var releaseCols = []ColumnDefinition{
	{title: "Name", flexFactor: 1},
	{title: "Namespace", flexFactor: 1},
	{title: "Revision", width: 10},
	{title: "Updated", width: 36},
	{title: "Status", flexFactor: 1},
	{title: "Chart", flexFactor: 1},
	{title: "App version", flexFactor: 1},
}

var historyCols = []ColumnDefinition{
	{title: "Revision", flexFactor: 1},
	{title: "Updated", width: 36},
	{title: "Status", flexFactor: 1},
	{title: "Chart", flexFactor: 1},
	{title: "App version", flexFactor: 1},
	{title: "Description", flexFactor: 1},
}

var releaseTableCache table.Model

func generateTables() (table.Model, table.Model) {
	t := table.New()
	h := table.New()
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
	h.SetStyles(s)
	t.KeyMap = k
	h.KeyMap = k
	return t, h
}

func (m Model) setTable(t *table.Model, cols []ColumnDefinition, targetWidth int, targetHeight int) tea.Cmd {
	var columns = make([]table.Column, len(cols))
	targetWidth = targetWidth - 2 // remove the border width
	remainingWidthAfterFixed := targetWidth
	totalFlex := 0
	for _, col := range cols {
		if col.flexFactor != 0 {
			totalFlex += col.flexFactor
		}
	}
	for i, col := range cols {
		if col.width != 0 {
			columns[i] = table.Column{
				Title: col.title,
				Width: col.width, // remove cell padding
			}
			remainingWidthAfterFixed = remainingWidthAfterFixed - col.width
		}
	}
	for i, col := range cols {
		if col.flexFactor != 0 {
			columns[i] = table.Column{
				Title: col.title,
				Width: int(remainingWidthAfterFixed*col.flexFactor/totalFlex) - 2, // -2 to remove the cell padding
			}
		}
	}
	// fill last column with the remaning width due to integer division
	lastCol := columns[len(columns)-1]
	totalColWidth := 0
	for _, col := range columns {
		totalColWidth += col.Width + 2 // count the cell padding
	}
	lastCol.Width = lastCol.Width + targetWidth - totalColWidth
	columns[len(columns)-1] = lastCol
	t.SetColumns(columns)
	t.SetHeight(targetHeight)
	t.SetWidth(targetWidth)
	return nil
}
