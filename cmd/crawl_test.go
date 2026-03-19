package cmd

import (
	"context"
	"strings"
	"testing"
)

func TestCrawlCmd_CancelWithURL_Rejected(t *testing.T) {
	cmd := &CrawlCmd{
		Target: "https://example.com",
		Cancel: true,
	}

	// handleURL should reject --cancel before touching the client (nil is fine)
	err := cmd.handleURL(context.Background(), nil, cmd.Target)
	if err == nil {
		t.Fatal("expected error when --cancel is used with a URL")
	}
	if got := err.Error(); got == "" {
		t.Error("error message should not be empty")
	}
}

func TestCrawlCmd_CancelWithURL_ErrorMessage(t *testing.T) {
	cmd := &CrawlCmd{
		Target: "https://example.com",
		Cancel: true,
	}

	err := cmd.handleURL(context.Background(), nil, cmd.Target)
	if err == nil {
		t.Fatal("expected error")
	}
	want := "--cancel requires a job ID"
	if got := err.Error(); !strings.Contains(got, want) {
		t.Errorf("error = %q, should contain %q", got, want)
	}
}
