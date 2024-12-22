package releases

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pidanou/helm-tui/helpers"
	"github.com/pidanou/helm-tui/styles"
)

func (m Model) View() string {
	var view string
	if m.installing {
		return m.installModel.View()
	}
	if m.upgrading {
		return m.upgradeModel.View()
	}
	if m.deleting {
		confirmMsg := "  No release selected. Press n to go back  "
		if m.releaseTable.SelectedRow() != nil {
			confirmMsg = fmt.Sprintf("  Delete release %s? y/n  ", m.releaseTable.SelectedRow()[0])
		}
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, styles.ActiveStyle.Border(styles.Border).Render(confirmMsg))
	}

	switch m.selectedView {
	case releasesView:
		tHeight := m.height - 2 - 1 // releaseTable padding + helper
		m.releaseTable.SetHeight(tHeight)
		view = m.renderReleasesTableView()
	default:
		view = m.renderReleaseDetail()
	}

	helperStyle := m.help.Styles.ShortSeparator
	helpView := m.help.View(m.keys[m.selectedView]) + helperStyle.Render(" • ") + m.help.View(helpers.CommonKeys)
	return view + "\n" + helpView
}

func (m Model) menuView() string {
	doc := strings.Builder{}

	var renderedTabs []string

	for i, t := range menuItem {
		var style lipgloss.Style
		isFirst, isActive := i == 0, i == int(m.selectedView)-1
		if isActive {
			style = styles.ActiveTabStyle
		} else {
			style = styles.InactiveTabStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		}
		style = style.Border(border)
		renderedTabs = append(renderedTabs, style.Render(t))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
	doc.WriteString(row + strings.Repeat("─", m.width-lipgloss.Width(row)-1) + styles.Border.TopRight)
	return doc.String()
}

func (m Model) renderReleaseDetail() string {
	header := m.renderReleasesTableView() + "\n" + m.menuView()
	remainingHeight := m.height - lipgloss.Height(header) + lipgloss.Height(m.menuView()) - 2 - 1 // releaseTable padding + helper
	var view string
	switch m.selectedView {
	case historyView:
		m.historyTable.SetHeight(remainingHeight - 2)
		view = header + "\n" + m.renderHistoryTableView()
	case notesView:
		m.notesVP.Height = remainingHeight - 4
		view = header + "\n" + m.renderNotesView()
	case metadataView:
		m.metadataVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderMetadataView()
	case hooksView:
		m.hooksVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderHooksView()
	case valuesView:
		m.valuesVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderValuesView()
	case manifestView:
		m.manifestVP.Height = remainingHeight - 4 // -4: 2*1 Padding + 2 borders
		view = header + "\n" + m.renderManifestView()
	}
	return view
}

func (m Model) renderReleasesTableView() string {
	var releasesTopBorder string
	tableView := m.releaseTable.View()
	var baseStyle lipgloss.Style
	releasesTopBorder = styles.GenerateTopBorderWithTitle(" Releases ", m.releaseTable.Width(), styles.Border, styles.InactiveStyle)
	baseStyle = styles.InactiveStyle.Border(styles.Border, false, true, true)
	tableView = baseStyle.Render(tableView)
	return lipgloss.JoinVertical(lipgloss.Top, releasesTopBorder, tableView)
}

func (m Model) renderHistoryTableView() string {
	tableView := m.historyTable.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Border(styles.Border).UnsetBorderTop()
	tableView = baseStyle.Render(tableView)
	return tableView
}

func (m Model) renderNotesView() string {
	view := m.notesVP.View()
	var baseStyle lipgloss.Style
	baseStyle = styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderMetadataView() string {
	view := m.metadataVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderHooksView() string {
	view := m.hooksVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderValuesView() string {
	view := m.valuesVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}

func (m Model) renderManifestView() string {
	view := m.manifestVP.View()
	baseStyle := styles.InactiveStyle.Padding(1, 2).Border(styles.Border, false, true, true)
	view = baseStyle.Render(view)
	return view
}
