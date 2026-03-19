package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout runs fn and returns what it wrote to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestReadOutput_SummaryUsesTarget(t *testing.T) {
	// Build content long enough to trigger structural summary (>2000 lines)
	var b strings.Builder
	b.WriteString("# Section One\n\n")
	for i := 0; i < 2100; i++ {
		b.WriteString("Line of content.\n")
	}
	content := b.String()

	cmd := &ReadCmd{URL: ""} // empty — URL came from -d body
	target := "https://from-data-body.com/docs"

	out := captureStdout(t, func() {
		cmd.output("/tmp/cache.md", target, content, true)
	})

	if !strings.Contains(out, "[ctx:summary]") {
		t.Fatal("expected structural summary for long document")
	}
	if !strings.Contains(out, target) {
		t.Errorf("summary should reference target URL %q, got:\n%s", target, out[:min(len(out), 200)])
	}
	if strings.Contains(out, "ctx read  -s") {
		t.Error("summary contains empty URL (double space), target was not passed through")
	}
}

func TestReadOutput_SummaryUsesExplicitURL(t *testing.T) {
	var b strings.Builder
	b.WriteString("# Heading\n\n")
	for i := 0; i < 2100; i++ {
		b.WriteString("Content line.\n")
	}
	content := b.String()

	target := "https://explicit.com/page"
	cmd := &ReadCmd{URL: target}

	out := captureStdout(t, func() {
		cmd.output("/tmp/cache.md", target, content, true)
	})

	if !strings.Contains(out, "ctx read "+target+" -s") {
		t.Errorf("summary should contain navigation hint with URL, got:\n%s", out[:min(len(out), 200)])
	}
}

func TestReadOutput_ShortDocNeverSummarizes(t *testing.T) {
	content := "# Hello\n\nShort content.\n"
	cmd := &ReadCmd{URL: ""}

	out := captureStdout(t, func() {
		cmd.output("/tmp/cache.md", "https://example.com", content, true)
	})

	if strings.Contains(out, "[ctx:summary]") {
		t.Error("short doc should not produce structural summary")
	}
	if out != content {
		t.Errorf("short doc should be printed as-is, got:\n%s", out)
	}
}
