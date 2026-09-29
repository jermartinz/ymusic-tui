package tui

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestNewModelStartsWithSearchFocused(t *testing.T) {
	model := NewModel()

	if !model.search.Focused() {
		t.Fatal("expected search input to start focused")
	}
	if model.activePane != resultsPane {
		t.Fatalf("expected results pane to start active, got %d", model.activePane)
	}
}

func TestWindowSizeUpdatesLayoutState(t *testing.T) {
	model := NewModel()

	updated, _ := model.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	got := updated.(Model)

	if got.width != 120 || got.height != 40 {
		t.Fatalf("expected 120x40, got %dx%d", got.width, got.height)
	}
}

func TestTabMovesFocusBetweenPanes(t *testing.T) {
	model := NewModel()
	model.search.Blur()

	updated, _ := model.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	got := updated.(Model)

	if got.activePane != queuePane {
		t.Fatalf("expected queue pane to be active, got %d", got.activePane)
	}
}

func TestViewContainsMainSections(t *testing.T) {
	model := NewModel()
	view := model.View().Content

	for _, section := range []string{"YMusic TUI", "Search results", "Queue", "Nothing playing"} {
		if !strings.Contains(view, section) {
			t.Errorf("expected view to contain %q", section)
		}
	}
}

func TestEmptyStateMessage(t *testing.T) {
	tests := []struct {
		name string
		pane pane
		want string
	}{
		{
			name: "results pane",
			pane: resultsPane,
			want: "Search for a song to see results.",
		},
		{
			name: "queue pane",
			pane: queuePane,
			want: "Your queue is empty.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := emptyStateMessage(tt.pane)
			if got != tt.want {
				t.Errorf("emptyStateMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
