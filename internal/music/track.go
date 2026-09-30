// Package music contains the domain types shared by search, playback, and the UI.
package music

import "time"

// Track is the provider-independent representation of playable music.
type Track struct {
	ID       string
	Title    string
	Artist   string
	URL      string
	Duration time.Duration
	Source   string
}
