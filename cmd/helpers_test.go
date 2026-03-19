package cmd

import "testing"

func TestEffectiveURL_CLIArg(t *testing.T) {
	got := effectiveURL("https://from-cli.com", []byte(`{"url":"https://from-body.com"}`))
	if got != "https://from-cli.com" {
		t.Errorf("got %q, want CLI arg to win", got)
	}
}

func TestEffectiveURL_FromBody(t *testing.T) {
	got := effectiveURL("", []byte(`{"url":"https://from-body.com"}`))
	if got != "https://from-body.com" {
		t.Errorf("got %q, want URL from body", got)
	}
}

func TestEffectiveURL_NilBody(t *testing.T) {
	got := effectiveURL("", nil)
	if got != "<unknown URL>" {
		t.Errorf("got %q, want fallback", got)
	}
}

func TestEffectiveURL_BodyWithoutURL(t *testing.T) {
	got := effectiveURL("", []byte(`{"prompt":"test"}`))
	if got != "<unknown URL>" {
		t.Errorf("got %q, want fallback when body has no url field", got)
	}
}
