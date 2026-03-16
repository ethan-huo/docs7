package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const defaultBaseURL = "https://context7.com"

type Client struct {
	BaseURL    string
	httpClient *http.Client
}

func NewClient() *Client {
	base := os.Getenv("CONTEXT7_BASE_URL")
	if base == "" {
		base = defaultBaseURL
	}
	return &Client{
		BaseURL:    base,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) SearchLibraries(name, query string) ([]Library, error) {
	params := url.Values{"libraryName": {name}}
	if query != "" {
		params.Set("query", query)
	}

	var resp librarySearchResponse
	if err := c.get("/api/v2/libs/search", params, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("%s", resp.Error)
	}
	return resp.Results, nil
}

func (c *Client) QueryDocs(libraryID, query string) (*DocsResponse, error) {
	params := url.Values{
		"libraryId": {libraryID},
		"query":     {query},
		"type":      {"json"},
	}

	var resp DocsResponse
	if err := c.get("/api/v2/context", params, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, fmt.Errorf("%s", resp.Error)
	}
	return &resp, nil
}

func (c *Client) get(path string, params url.Values, out any) error {
	u := c.BaseURL + path + "?" + params.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return json.Unmarshal(body, out)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("X-Context7-Source", "cli")
	req.Header.Set("X-Context7-Client-IDE", "docs7")
	req.Header.Set("X-Context7-Transport", "cli")

	// GetValidToken handles env var, file load, expiry check, and refresh
	if token, _ := GetValidToken(c.BaseURL); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}
