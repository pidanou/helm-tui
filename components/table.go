package components

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
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
