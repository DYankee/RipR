package styles

import (
	"charm.land/lipgloss/v2"
)

// Colors
const (
	PrimaryColor    = "205"
	SecondaryColor  = "240"
	TertiaryColor   = "344"
	BackgroundColor = "#1a1a1a"
	AccentColor     = "86"
)

// Base styles
var (
	FocusedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(BackgroundColor)).
			Foreground(lipgloss.Color(PrimaryColor)).
			Bold(true)

	BlurStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(BackgroundColor)).
			Foreground(lipgloss.Color(SecondaryColor))

	NoStyle      = lipgloss.NewStyle()
	HelpStyle    = BlurStyle
	CursorStyle  = FocusedStyle
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	LoadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor))
	CursorMode   = lipgloss.NewStyle().Foreground(lipgloss.Color(TertiaryColor))

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = BlurStyle.Render("[ Submit ]")

	// Layout styles
	WindowStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(BackgroundColor)).
			Padding(1, 2).
			Margin(0)

	TitleStyle = FocusedStyle.
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(lipgloss.Color(PrimaryColor)).
			PaddingBottom(1).
			MarginBottom(1)

	SectionStyle = FocusedStyle.
			Bold(true).
			MarginBottom(1)

	DebugStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(TertiaryColor)).
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(SecondaryColor)).
			PaddingLeft(2).
			MarginTop(1)

	InputStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(BackgroundColor)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(SecondaryColor))

	InputFocusedStyle = InputStyle.
				BorderForeground(lipgloss.Color(PrimaryColor))
)

// Helper to create a box for content
func Box(content string, title string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(SecondaryColor)).
		Padding(1, 2)

	return style.Render(content)
}
