package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateCredentials_RefusesOnCorruptFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "ctx")
	os.MkdirAll(dir, 0o700)

	corrupt := []byte("{{{{ not valid yaml ::::")
	path := filepath.Join(dir, "credentials.yaml")
	os.WriteFile(path, corrupt, 0o600)

	err := UpdateCredentials(func(c *Credentials) {
		c.Cloudflare.APIToken = "should-never-be-written"
	})
	if err == nil {
		t.Fatal("expected error when credentials file is corrupt, got nil")
	}

	// Verify file was NOT overwritten
	after, _ := os.ReadFile(path)
	if string(after) != string(corrupt) {
		t.Errorf("corrupt file was overwritten:\nbefore: %s\nafter:  %s", corrupt, after)
	}
}

func TestUpdateCredentials_RefusesOnPermissionError(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	dir := filepath.Join(tmp, ".config", "ctx")
	os.MkdirAll(dir, 0o700)

	path := filepath.Join(dir, "credentials.yaml")
	os.WriteFile(path, []byte("cloudflare:\n  api_token: keep-me\n"), 0o600)

	// Remove read permission
	os.Chmod(path, 0o000)
	t.Cleanup(func() { os.Chmod(path, 0o600) })

	err := UpdateCredentials(func(c *Credentials) {
		c.Cloudflare.APIToken = "should-never-be-written"
	})
	if err == nil {
		t.Fatal("expected error when credentials file is unreadable, got nil")
	}
}

func TestUpdateCredentials_WorksOnMissingFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// No credentials.yaml exists — should create new one via migrateOldCredentials
	err := UpdateCredentials(func(c *Credentials) {
		c.Cloudflare.APIToken = "new-token"
	})
	if err != nil {
		t.Fatalf("unexpected error for missing file: %v", err)
	}

	// Verify file was created
	path := filepath.Join(tmp, ".config", "ctx", "credentials.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("credentials file not created: %v", err)
	}
	if len(data) == 0 {
		t.Error("credentials file is empty")
	}
}
