package tui

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
	wideLayoutAt  = 100
)

type pane int

const (
	resultsPane pane = iota
	queuePane
)

// Model contains all state required to update and render the terminal UI.
type Model struct {
	search     textinput.Model
	activePane pane
	width      int
	height     int
}

// NewModel creates the initial application state.
func NewModel() Model {
	search := textinput.New()
	search.Prompt = "> "
	search.Placeholder = "Search for music"
	search.CharLimit = 120
	search.SetWidth(defaultWidth - 8)
	search.SetVirtualCursor(true)
	search.Focus()

	return Model{
		search:     search,
		activePane: resultsPane,
		width:      defaultWidth,
		height:     defaultHeight,
	}
}

// Init starts the text input cursor.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update applies terminal events to the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.search.SetWidth(max(msg.Width-8, 12))

	case tea.KeyPressMsg:
		key := msg.String()

		if key == "ctrl+c" {
			return m, tea.Quit
		}

		if m.search.Focused() {
			switch key {
			case "esc", "enter":
				m.search.Blur()
				return m, nil
			case "tab":
				m.search.Blur()
				m.activePane = resultsPane
				return m, nil
			}
		} else {
			switch key {
			case "q":
				return m, tea.Quit
			case "/":
				return m, m.search.Focus()
			case "tab":
				m.activePane = m.activePane.next()
				return m, nil
			}
		}
	}

	if m.search.Focused() {
		var cmd tea.Cmd
		m.search, cmd = m.search.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the complete terminal screen from the current model.
func (m Model) View() tea.View {
	width := max(m.width, 30)
	height := max(m.height, 16)
	contentWidth := width - 2
	bodyHeight := max(height-12, 6)

	header := titleStyle.Render("YMusic TUI") + "  " + mutedStyle.Render("public audio · provider: auto")
	searchBox := searchStyle.Width(max(contentWidth-4, 16)).Render(m.search.View())
	body := m.renderBody(contentWidth, bodyHeight)
	player := playerStyle.Width(max(contentWidth-4, 16)).Render(
		mutedStyle.Render("Nothing playing") + "  " + progressStyle.Render("────────────────────") + "  00:00 / 00:00",
	)
	help := mutedStyle.Render("/ search  •  tab switch pane  •  j/k navigate  •  ? help  •  q quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		searchBox,
		body,
		player,
		help,
	)

	view := tea.NewView(appStyle.Width(contentWidth).Render(content))
	view.AltScreen = true
	view.WindowTitle = "YMusic TUI"
	return view
}

func (m Model) renderBody(width, height int) string {
	if width >= wideLayoutAt {
		gap := 1
		leftWidth := (width - gap) / 2
		rightWidth := width - gap - leftWidth

		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.renderPanel("Search results", emptyStateMessage(resultsPane), leftWidth, height, m.activePane == resultsPane),
			" ",
			m.renderPanel("Queue", emptyStateMessage(queuePane), rightWidth, height, m.activePane == queuePane),
		)
	}

	topHeight := height / 2
	bottomHeight := height - topHeight
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderPanel("Search results", emptyStateMessage(resultsPane), width, topHeight, m.activePane == resultsPane),
		m.renderPanel("Queue", emptyStateMessage(queuePane), width, bottomHeight, m.activePane == queuePane),
	)
}

func (m Model) renderPanel(title, content string, width, height int, focused bool) string {
	borderColor := mutedColor
	if focused {
		borderColor = accentColor
	}

	style := panelStyle.
		BorderForeground(borderColor).
		Width(max(width-4, 12)).
		Height(max(height-2, 1))

	return style.Render(sectionTitleStyle.Render(title) + "\n\n" + mutedStyle.Render(content))
}

func (p pane) next() pane {
	if p == resultsPane {
		return queuePane
	}
	return resultsPane
}

// emptyStateMessage explains the next action available in an empty pane.
func emptyStateMessage(p pane) string {
	switch p {
	case resultsPane:
		return "Search for a song to see results."
	case queuePane:
		return "Your queue is empty."
	default:
		return "Nothing here yet."
	}
}
