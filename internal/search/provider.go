// Package search defines provider-independent music search behavior.
package search

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jermartinz/ymusic-tui/internal/music"
)

const (
	// MaxQueryLength keeps provider requests within a predictable size.
	MaxQueryLength = 120
)

var (
	// ErrEmptyQuery is returned when a query contains no searchable text.
	ErrEmptyQuery = errors.New("search query cannot be empty")
	// ErrQueryTooLong is returned when a query exceeds MaxQueryLength.
	ErrQueryTooLong = errors.New("search query is too long")
)

// Kind identifies a search provider or the automatic selection mode.
type Kind string

const (
	KindAuto    Kind = "auto"
	KindYouTube Kind = "youtube"
	KindYTDLP   Kind = "ytdlp"
)

// Provider searches for music and returns normalized tracks.
type Provider interface {
	Kind() Kind
	Search(ctx context.Context, query string, limit int) ([]music.Track, error)
}

// ParseKind converts user input into a supported provider kind.
func ParseKind(value string) (Kind, error) {
	kind := Kind(strings.ToLower(strings.TrimSpace(value)))
	if kind == "" {
		return KindAuto, nil
	}

	switch kind {
	case KindAuto, KindYouTube, KindYTDLP:
		return kind, nil
	default:
		return "", fmt.Errorf("unsupported search provider %q", value)
	}
}

// ResolveKind chooses the concrete provider used for this application run.
func ResolveKind(requested Kind, youtubeAPIKey string) (Kind, error) {
	switch requested {
	case KindYouTube, KindYTDLP:
		return requested, nil
	case KindAuto:
		if strings.TrimSpace(youtubeAPIKey) != "" {
			return KindYouTube, nil
		}
		return KindYTDLP, nil
	default:
		return "", fmt.Errorf("unsupported search provider %q", requested)
	}
}

// ValidateQuery checks whether a query is safe to send to a provider.
func ValidateQuery(query string) error {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return ErrEmptyQuery
	}
	if len([]rune(trimmed)) > MaxQueryLength {
		return ErrQueryTooLong
	}
	return nil
}
