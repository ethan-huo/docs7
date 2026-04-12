package cmd

import (
	"strings"
	"testing"
)

func TestRenderDocSources_LinksTitlesAndOmitsURLList(t *testing.T) {
	sources := []docSource{
		{
			title:        "Ghostty Documentation - AppleScript (macOS)",
			url:          "https://ghostty.org/docs/features/applescript",
			descriptions: []string{"Split the current terminal horizontally.", "Send text to the new pane."},
		},
		{
			title:        "Unknown",
			url:          "https://ghostty.org/docs/vt/control/tab",
			descriptions: []string{"Documentation for the Tab control character."},
		},
	}

	var out strings.Builder
	renderDocSources(&out, sources)

	got := out.String()
	want := strings.Join([]string{
		"## 1. [Ghostty Documentation - AppleScript (macOS)](https://ghostty.org/docs/features/applescript)",
		"- Split the current terminal horizontally.",
		"- Send text to the new pane.",
		"",
		"## 2. https://ghostty.org/docs/vt/control/tab",
		"- Documentation for the Tab control character.",
		"",
		"---",
		"Use `ctx read <url>` to fetch full documents.",
		"",
	}, "\n")

	if got != want {
		t.Fatalf("renderDocSources() =\n%q\nwant\n%q", got, want)
	}
}

func TestDocSourceHeading_LinksTitlesAndFallsBackToBareURL(t *testing.T) {
	tests := []struct {
		name string
		src  docSource
		want string
	}{
		{
			name: "blank",
			src:  docSource{title: " ", url: "https://example.com/blank"},
			want: "https://example.com/blank",
		},
		{
			name: "unknown",
			src:  docSource{title: "Unknown", url: "https://example.com/unknown"},
			want: "https://example.com/unknown",
		},
		{
			name: "known",
			src:  docSource{title: "Guide", url: "https://example.com/guide"},
			want: "[Guide](https://example.com/guide)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := docSourceHeading(tt.src)
			if got != tt.want {
				t.Fatalf("docSourceHeading() = %q, want %q", got, tt.want)
			}
		})
	}
}
