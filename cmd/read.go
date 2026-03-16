package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/anthropics/docs7/api"
)

type ReadCmd struct {
	URL string `arg:"" help:"URL to read (github://owner/repo/path or https://...)"`
}

func (c *ReadCmd) Run(_ *api.Client) error {
	url := c.URL

	// github://owner/repo/path -> gh API
	if strings.HasPrefix(url, "github://") {
		return readGitHub(strings.TrimPrefix(url, "github://"))
	}

	// https://github.com/owner/repo/blob/branch/path -> gh API
	if strings.Contains(url, "github.com") {
		if path, ok := parseGitHubBlobURL(url); ok {
			return readGitHub(path)
		}
	}

	// Everything else -> HTTP with markdown content negotiation
	return readHTTP(url)
}

// readGitHub fetches file content via GitHub API.
// path format: owner/repo/path/to/file.md
func readGitHub(path string) error {
	parts := strings.SplitN(path, "/", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid github path: %s (expected owner/repo/path)", path)
	}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", parts[0], parts[1], parts[2])

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// Use gh's token if available
	if token := ghToken(); token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(
			strings.ReplaceAll(result.Content, "\n", ""),
		)
		if err != nil {
			return err
		}
		fmt.Print(string(decoded))
		return nil
	}

	fmt.Print(result.Content)
	return nil
}

// readHTTP fetches URL with markdown content negotiation.
func readHTTP(url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Prefer markdown, fall back to HTML
	req.Header.Set("Accept",
		"text/markdown;q=1.0, text/x-markdown;q=0.9, text/plain;q=0.8, text/html;q=0.7, */*;q=0.1")
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko)")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	ct := resp.Header.Get("Content-Type")
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024)) // 5MB cap
	if err != nil {
		return err
	}

	if isMarkdownContentType(ct) || isPlainText(ct) {
		fmt.Print(string(body))
		return nil
	}

	// HTML response — fallback to Jina Reader for markdown conversion
	return readJina(url)
}

func readJina(originalURL string) error {
	req, err := http.NewRequest("GET", "https://r.jina.ai/"+originalURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/markdown")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Jina Reader %d for %s", resp.StatusCode, originalURL)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5*1024*1024))
	if err != nil {
		return err
	}
	fmt.Print(string(body))
	return nil
}

func isMarkdownContentType(ct string) bool {
	return strings.Contains(ct, "text/markdown") || strings.Contains(ct, "text/x-markdown")
}

func isPlainText(ct string) bool {
	return strings.Contains(ct, "text/plain")
}

// parseGitHubBlobURL extracts owner/repo/path from a GitHub blob URL.
func parseGitHubBlobURL(rawURL string) (string, bool) {
	trimmed := strings.TrimPrefix(rawURL, "https://github.com/")
	trimmed = strings.TrimPrefix(trimmed, "http://github.com/")
	parts := strings.SplitN(trimmed, "/blob/", 2)
	if len(parts) != 2 {
		return "", false
	}
	repo := parts[0]
	rest := parts[1]
	// Strip branch prefix
	if idx := strings.IndexByte(rest, '/'); idx >= 0 {
		rest = rest[idx+1:]
	} else {
		return "", false
	}
	// Strip query/fragment
	if idx := strings.IndexByte(rest, '?'); idx >= 0 {
		rest = rest[:idx]
	}
	if idx := strings.IndexByte(rest, '#'); idx >= 0 {
		rest = rest[:idx]
	}
	return repo + "/" + rest, true
}

func ghToken() string {
	if t := os.Getenv("GITHUB_TOKEN"); t != "" {
		return t
	}
	if t := os.Getenv("GH_TOKEN"); t != "" {
		return t
	}
	// Fallback: gh auth token
	out, err := exec.Command("gh", "auth", "token").Output()
	if err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

var httpClient = &http.Client{Timeout: 30 * time.Second}
