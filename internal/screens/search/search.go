package search

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"

	// Internal
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/DYankee/RRipper/internal/nav"
	"github.com/DYankee/RRipper/internal/styles"
)

type SearchResultMsg struct {
	releases []mb.Release
	err      error
}

type Model struct {
	client *mb.MusicBrainz
	inputs []textinput.Model

	focused    int
	loading    bool
	cursorMode cursor.Mode
	err        error
}

func New(client *mb.MusicBrainz) Model {
	m := Model{
		client: client,
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CharLimit = 100
		switch i {
		case 0:
			t.SetWidth(20)
			t.Placeholder = "Artist"
			t.Focus()
		case 1:
			t.SetWidth(20)
			t.Placeholder = "Album"
		}
		m.inputs[i] = t
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// SUBMIT LOGIC: Instead of tea.Quit, send the message
			if s == "enter" && m.focused == len(m.inputs) {
				m.loading = true
				return m, m.doSearch()
			}

			// Cycle focus
			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused > len(m.inputs) {
				m.focused = 0
			} else if m.focused < 0 {
				m.focused = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focused {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)
		}
	case SearchResultMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		return m, func() tea.Msg {
			return nav.ToResults{Releases: msg.releases}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var b strings.Builder

	// Title
	b.WriteString(styles.TitleStyle.Render("Search MusicBrainz"))
	b.WriteString("\n\n")

	// Instructions
	b.WriteString(styles.BlurStyle.Render("Find vinyl releases by artist and album"))
	b.WriteString("\n\n")

	// Input fields
	labels := []string{"Artist:", "Album:"}
	for i := range m.inputs {
		prefix := "  "
		if i == m.focused && !m.loading {
			prefix = styles.FocusedStyle.Render("> ")
		}

		label := styles.BlurStyle.Render(labels[i])
		input := styles.BlurStyle.Render(m.inputs[i].View())
		b.WriteString(styles.InputFocusedStyle.Render(fmt.Sprintf("%s%-10s%s", prefix, label, input)))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Submit button
	submitIdx := len(m.inputs)
	buttonPrefix := "  "
	if m.focused == submitIdx {
		buttonPrefix = styles.FocusedStyle.Render("> ")
	}

	if m.loading {
		b.WriteString(styles.LoadingStyle.Render(buttonPrefix + "[ Searching... ]"))
	} else {
		button := styles.BlurredButton
		if m.focused == submitIdx {
			button = styles.FocusedButton
		}
		b.WriteString(buttonPrefix)
		b.WriteString(button)
	}

	// Error message
	if m.err != nil {
		b.WriteString("\n\n")
		b.WriteString(styles.ErrorStyle.Render("Error: " + m.err.Error()))
	}

	// Help text
	b.WriteString("\n\n")
	b.WriteString(styles.HelpStyle.Render("[↑/↓] Navigate  [enter] Search  [esc] Back"))

	// v := tea.NewView(styles.WindowStyle.Render(b.String()))
	return tea.NewView(b.String())
}

func (m Model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m Model) doSearch() tea.Cmd {
	artist := m.inputs[0].Value()
	release := m.inputs[1].Value()
	format := "12vinyl"
	return func() tea.Msg {
		releases, err := m.client.SearchRelease(artist, release, format)
		return SearchResultMsg{releases: releases, err: err}
	}
}
