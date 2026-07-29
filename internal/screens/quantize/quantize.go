package quantize

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	mb "github.com/DYankee/RRipper/internal/musicbrainz"
	"github.com/DYankee/RRipper/internal/nav"
	"github.com/DYankee/RRipper/internal/styles"
)

type track struct {
	title      string
	originalMs int
	input      textinput.Model
}

type Model struct {
	release *mb.Release
	tracks  []track
	focused int
}

func New(release *mb.Release, w, h int) Model {
	tracks := make([]track, len(release.Tracks))
	for i, t := range release.Tracks {
		inp := newTrackInput(formatMs(quantize(t.Length)))
		tracks[i] = track{
			title:      t.Title,
			originalMs: t.Length,
			input:      inp,
		}
	}
	if len(tracks) > 0 {
		tracks[0].input.Focus()
	}
	return Model{release: release, tracks: tracks, focused: 0}
}

// newTrackInput creates a textinput with focused/blurred styles pre-configured.
// Since v2 stores both states in a Styles struct, we set them once here and
// never need to touch styles again in moveFocus — just call Focus()/Blur().
func newTrackInput(value string) textinput.Model {
	inp := textinput.New()
	inp.SetWidth(8)
	inp.SetValue(value)

	s := textinput.DefaultDarkStyles()
	s.Focused.Prompt = styles.FocusedStyle
	s.Focused.Text = styles.FocusedStyle
	s.Blurred.Prompt = styles.NoStyle
	s.Blurred.Text = styles.NoStyle
	inp.SetStyles(s)

	return inp
}

func (m Model) Init() tea.Cmd { return textinput.Blink }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			return m.moveFocus(1)
		case "shift+tab", "up":
			return m.moveFocus(-1)
		case "esc":
			return m, func() tea.Msg { return nav.GoBack{} }
		}
	}

	var cmd tea.Cmd
	m.tracks[m.focused].input, cmd = m.tracks[m.focused].input.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	var b strings.Builder

	b.WriteString(styles.FocusedStyle.Render("Quantized Track Lengths"))
	b.WriteString("\n")
	b.WriteString(
		styles.BlurStyle.Render(m.release.Title + " — " + m.release.Artist),
	)
	b.WriteString("\n\n")
	b.WriteString(styles.BlurStyle.Render(
		fmt.Sprintf("  %-4s  %-40s  %-8s  %-8s\n", "#", "Title", "Original", "Adjusted"),
	))

	for i, t := range m.tracks {
		cursor := "  "
		if i == m.focused {
			cursor = styles.FocusedStyle.Render("> ")
		}
		b.WriteString(fmt.Sprintf("%s%-4d  %-40s  %-8s  %s\n",
			cursor,
			i+1,
			t.title,
			formatMs(t.originalMs),
			t.input.View(),
		))
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("[↑/↓] Navigate  [esc] Back"))

	return tea.NewView(styles.WindowStyle.Render(b.String()))
}

// moveFocus is now just Focus/Blur — no style swapping needed since both
// states were baked into the input's Styles struct at creation.
func (m Model) moveFocus(dir int) (Model, tea.Cmd) {
	m.tracks[m.focused].input.Blur()
	m.focused = (m.focused + dir + len(m.tracks)) % len(m.tracks)
	m.tracks[m.focused].input.Focus()
	return m, nil
}

// quantize is a placeholder — fill in your actual Audacity quantization logic.
func quantize(ms int) int {
	return ms
}

func formatMs(ms int) string {
	s := ms / 1000
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}
