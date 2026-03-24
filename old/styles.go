package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	defaultBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder())

	focusedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205"))

	blurredStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))

	cursorStyle = focusedStyle
	noStyle     = lipgloss.NewStyle()

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf(
		"[ %s ]",
		blurredStyle.Render("Submit"),
	)

	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63"))
)
