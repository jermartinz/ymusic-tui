package tui

import "charm.land/lipgloss/v2"

var (
	accentColor = lipgloss.Color("#D75FFF")
	mutedColor  = lipgloss.Color("#6C7086")
	textColor   = lipgloss.Color("#CDD6F4")

	appStyle = lipgloss.NewStyle().
			Foreground(textColor).
			Padding(1)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(textColor)

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	searchStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	playerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(0, 1)

	progressStyle = lipgloss.NewStyle().Foreground(accentColor)
)
