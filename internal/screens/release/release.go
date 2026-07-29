package release

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	// Internal
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/DYankee/RRipper/internal/nav"
	"github.com/DYankee/RRipper/internal/styles"
)

type Model struct {
	release  *mb.Release
	viewport viewport.Model
}

func New(r *mb.Release, w int, h int) Model {
	vp := viewport.New(
		viewport.WithWidth(w),
		viewport.WithHeight(h-4),
	)
	vp.SetContent(renderRelease(r))
	return Model{release: r, viewport: vp}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// user confirmed this album — move to quantize
			return m, func() tea.Msg {
				return nav.ToQuantize{Release: m.release}
			}
		case "esc":
			return m, func() tea.Msg {
				return nav.ToResults{}
			}
		}
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	help := styles.HelpStyle.Render(
		"[↑/↓] Scroll  [enter] Confirm  [esc] Back",
	)
	v := tea.NewView(styles.WindowStyle.Render(m.viewport.View() + "\n\n" + help))
	return v
}

func renderRelease(r *mb.Release) string {
	var b strings.Builder

	b.WriteString(styles.FocusedStyle.Render(r.Title))
	b.WriteString("\n")
	b.WriteString(styles.BlurStyle.Render(r.Artist + " • " + r.Date))
	b.WriteString("\n\n")
	b.WriteString(styles.FocusedStyle.Render("Tracks\n"))

	for _, t := range r.Tracks {
		b.WriteString(fmt.Sprintf(
			"  %2d.  %-40s  %s\n",
			t.Position,
			t.Title,
			formatMs(t.Length),
		))
	}

	return b.String()
}

func formatMs(ms int) string {
	s := ms / 1000
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
