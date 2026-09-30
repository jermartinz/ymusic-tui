package search

import (
	"errors"
	"strings"
	"testing"
)

func TestParseKind(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    Kind
		wantErr bool
	}{
		{name: "empty means auto", value: "", want: KindAuto},
		{name: "normalizes case and spaces", value: " YouTube ", want: KindYouTube},
		{name: "accepts yt-dlp provider", value: "ytdlp", want: KindYTDLP},
		{name: "rejects unknown provider", value: "spotify", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseKind(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseKind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseKind() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveKind(t *testing.T) {
	tests := []struct {
		name      string
		requested Kind
		apiKey    string
		want      Kind
		wantErr   bool
	}{
		{name: "auto uses YouTube with API key", requested: KindAuto, apiKey: "secret", want: KindYouTube},
		{name: "auto uses yt-dlp without API key", requested: KindAuto, want: KindYTDLP},
		{name: "explicit YouTube wins without key", requested: KindYouTube, want: KindYouTube},
		{name: "explicit yt-dlp wins with key", requested: KindYTDLP, apiKey: "secret", want: KindYTDLP},
		{name: "unknown provider fails", requested: Kind("unknown"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveKind(tt.requested, tt.apiKey)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ResolveKind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ResolveKind() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestValidateQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr error
	}{
		{
			name:  "valid query",
			query: "Daft Punk",
		},
		{
			name:    "empty query",
			query:   "  ",
			wantErr: ErrEmptyQuery,
		},
		{
			name:    "query over maximum length",
			query:   strings.Repeat("a", MaxQueryLength+1),
			wantErr: ErrQueryTooLong,
		},
		{
			name:  "query at maximum length",
			query: strings.Repeat("a", MaxQueryLength),
		},
		{
			name:  "ignores surrounding spaces at maximum length",
			query: " " + strings.Repeat("a", MaxQueryLength) + " ",
		},
		{
			name:  "counts Unicode runes at maximum length",
			query: strings.Repeat("á", MaxQueryLength),
		},
		{
			name:    "rejects Unicode query over maximum length",
			query:   strings.Repeat("á", MaxQueryLength+1),
			wantErr: ErrQueryTooLong,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateQuery(tt.query)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf(
					"ValidateQuery() error = %v, want %v",
					err,
					tt.wantErr,
				)
			}
		})
	}

}
