package main

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	switch m.currentView {
	case viewSearch:
		return m.searchView()
	case viewLoading:
		return m.loadingView()
	case viewSearchResult:
		return m.searchResultView()
	case viewReleaseResult:
		return m.releaseView()
	case viewExporting:
		return m.exportingView()
	default:
		return ""
	}
}

func (m *model) header() string {
	style := borderStyle.Copy().
		Width(m.width - 2).
		Align(lipgloss.Center)

	head := lipgloss.JoinHorizontal(
		lipgloss.Center,
		"| Current view: "+m.currentView.String()+" |",
		"| Index: "+strconv.Itoa(m.focusIndex)+" |",
	)

	switch m.currentView {
	case viewSearchResult:
		head = lipgloss.JoinVertical(
			lipgloss.Center,
			head,
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				"| Artist: "+m.inputs[0].Value()+" |",
				"| Release: "+m.inputs[1].Value()+" |",
			),
		)
	case viewReleaseResult:
		head = lipgloss.JoinVertical(
			lipgloss.Center,
			head,
			"| Release: "+m.releaseData.Title+
				" | "+m.releaseData.Disambiguation+" |",
			"| Cursor: "+strconv.Itoa(m.releaseDataTable.Cursor())+" |",
		)
	}

	return style.Render(head)
}

func (m *model) bodyStyle(extraPadTop int) lipgloss.Style {
	headerH := lipgloss.Height(m.header())
	return borderStyle.Copy().
		Width(m.width - 2).
		PaddingTop((m.height / 2) - (headerH + extraPadTop)).
		PaddingBottom((m.height / 2) - (headerH + 2)).
		AlignHorizontal(lipgloss.Center)
}

func (m *model) loadingView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		"Loading...",
	)
}

func (m *model) searchView() string {
	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		m.inputs[0].View(),
		m.inputs[1].View(),
		m.inputs[2].View(),
		*button,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.header(),
			m.bodyStyle(0).Render(body),
		),
	)
}

func (m *model) searchResultView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.header(),
			m.bodyStyle(20).Render(m.searchResTable.View()),
		),
	)
}

func (m *model) releaseView() string {
	button := &blurredButton
	if m.focusIndex == 1 {
		button = &focusedButton
	}

	body := lipgloss.JoinVertical(
		lipgloss.Center,
		m.releaseDataTable.View(),
		*button,
	)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.header(),
			m.bodyStyle(20).Render(body),
		),
	)
}

func (m *model) exportingView() string {
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		"Exporting...",
	)
}
