package cfrender

import (
	"os"
	"path/filepath"
	"testing"
)

// ===== ParseBody =====

func TestParseBody_Empty(t *testing.T) {
	f := DataFlag{Data: ""}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		t.Errorf("empty -d should return nil, got %q", body)
	}
}

func TestParseBody_InlineJSON5(t *testing.T) {
	f := DataFlag{Data: `{url: "https://example.com", timeout: 5}`}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatal(err)
	}

	// Should produce valid standard JSON
	want := `{"timeout":5,"url":"https://example.com"}`
	if string(body) != want {
		t.Errorf("body = %s, want %s", body, want)
	}
}

func TestParseBody_InlineStandardJSON(t *testing.T) {
	f := DataFlag{Data: `{"url":"https://example.com"}`}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatal(err)
	}
	if body == nil {
		t.Fatal("expected non-nil body")
	}
}

func TestParseBody_AtFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "body.json5")
	os.WriteFile(fp, []byte(`{url: "https://example.com"}`), 0o644)

	f := DataFlag{Data: "@" + fp}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatal(err)
	}
	if body == nil {
		t.Fatal("expected non-nil body from @file")
	}
}

func TestParseBody_AtFileMissing(t *testing.T) {
	f := DataFlag{Data: "@/nonexistent/file.json5"}
	_, err := f.ParseBody()
	if err == nil {
		t.Error("expected error for missing @file")
	}
}

func TestParseBody_InvalidJSON5(t *testing.T) {
	f := DataFlag{Data: `{invalid json5 !!!`}
	_, err := f.ParseBody()
	if err == nil {
		t.Error("expected error for invalid JSON5")
	}
}

func TestParseBody_TrailingComma(t *testing.T) {
	// JSON5 allows trailing commas
	f := DataFlag{Data: `{url: "https://example.com",}`}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatalf("JSON5 trailing comma should be valid: %v", err)
	}
	if body == nil {
		t.Fatal("expected non-nil body")
	}
}

func TestParseBody_SingleQuotes(t *testing.T) {
	// JSON5 allows single quotes
	f := DataFlag{Data: `{url: 'https://example.com'}`}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatalf("JSON5 single quotes should be valid: %v", err)
	}
	if body == nil {
		t.Fatal("expected non-nil body")
	}
}

func TestParseBody_Comments(t *testing.T) {
	// JSON5 allows comments
	f := DataFlag{Data: `{
		// this is a comment
		url: "https://example.com"
	}`}
	body, err := f.ParseBody()
	if err != nil {
		t.Fatalf("JSON5 comments should be valid: %v", err)
	}
	if body == nil {
		t.Fatal("expected non-nil body")
	}
}

// ===== HasData =====

func TestHasData(t *testing.T) {
	if (&DataFlag{}).HasData() {
		t.Error("empty DataFlag.HasData() should be false")
	}
	if !(&DataFlag{Data: "something"}).HasData() {
		t.Error("non-empty DataFlag.HasData() should be true")
	}
}

// ===== ResolveValue =====

func TestResolveValue_Plain(t *testing.T) {
	got, err := ResolveValue("Bearer token123")
	if err != nil {
		t.Fatal(err)
	}
	if got != "Bearer token123" {
		t.Errorf("got %q, want %q", got, "Bearer token123")
	}
}

func TestResolveValue_AtFile(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "token.txt")
	os.WriteFile(fp, []byte("  secret-value  \n"), 0o644)

	got, err := ResolveValue("@" + fp)
	if err != nil {
		t.Fatal(err)
	}
	if got != "secret-value" {
		t.Errorf("got %q, want %q (trimmed)", got, "secret-value")
	}
}

func TestResolveValue_AtFileMissing(t *testing.T) {
	_, err := ResolveValue("@/nonexistent/token.txt")
	if err == nil {
		t.Error("expected error for missing @file")
	}
}
