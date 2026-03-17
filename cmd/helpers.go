package cmd

import "encoding/json"

func effectiveURL(cliURL string, body []byte) string {
	if cliURL != "" {
		return cliURL
	}
	var m map[string]any
	if json.Unmarshal(body, &m) == nil {
		if u, ok := m["url"].(string); ok {
			return u
		}
	}
	return "<unknown URL>"
}
