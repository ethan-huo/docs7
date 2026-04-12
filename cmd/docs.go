package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ethan-huo/ctx/api"
)

type DocsCmd struct {
	Name  string `arg:"" help:"Library name or ID (e.g. 'mlx-swift' or '/ml-explore/mlx-swift')"`
	Query string `arg:"" help:"What to search for in the documentation"`
}

func (c *DocsCmd) Run(client *api.Client) error {
	libraryID := c.Name

	if !strings.HasPrefix(libraryID, "/") {
		libs, err := client.SearchLibraries(c.Name, c.Query)
		if err != nil {
			return err
		}
		if len(libs) == 0 {
			fmt.Printf("No libraries found for %q. Try: ctx search %s\n", c.Name, c.Name)
			return nil
		}
		libraryID = libs[0].ID
		fmt.Printf("# %s\n\n", libraryID)
	}

	resp, err := client.QueryDocs(libraryID, c.Query)
	if err != nil {
		return err
	}

	sources := extractSources(resp)
	if len(sources) == 0 {
		fmt.Printf("No documentation matched query %q in %s. Try a broader query.\n", c.Query, libraryID)
		return nil
	}

	renderDocSources(os.Stdout, sources)

	return nil
}

type docSource struct {
	title        string
	descriptions []string
	url          string
}

func extractSources(resp *api.DocsResponse) []docSource {
	seen := map[string]int{}
	var sources []docSource

	add := func(rawURL, pageTitle, desc string) {
		if !isRealURL(rawURL) {
			return
		}
		url := normalizeURL(rawURL)
		if idx, ok := seen[url]; ok {
			if desc != "" {
				sources[idx].descriptions = append(sources[idx].descriptions, desc)
			}
			return
		}
		src := docSource{
			title: pageTitle,
			url:   url,
		}
		if desc != "" {
			src.descriptions = []string{desc}
		}
		seen[url] = len(sources)
		sources = append(sources, src)
	}

	for _, s := range resp.CodeSnippets {
		desc := s.CodeDescription
		if desc == "" {
			desc = s.CodeTitle
		}
		add(s.CodeID, s.PageTitle, desc)
	}
	for _, s := range resp.InfoSnippets {
		title := s.Breadcrumb
		if title == "" {
			title = "Info"
		}
		add(s.URL, title, "")
	}
	return sources
}

func isRealURL(u string) bool {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return false
	}
	if strings.Contains(u, "context7_") || strings.Contains(u, "context7.com/") {
		return false
	}
	return true
}

// normalizeURL converts GitHub blob URLs to github:// scheme (preserving ref),
// keeps everything else as-is.
func normalizeURL(rawURL string) string {
	return canonicalizeURL(rawURL)
}

func renderDocSources(w io.Writer, sources []docSource) {
	for i, src := range sources {
		if i > 0 {
			fmt.Fprintln(w)
		}
		heading := docSourceHeading(src)
		fmt.Fprintf(w, "## %d. %s\n", i+1, heading)
		for j, desc := range src.descriptions {
			if j >= 3 {
				break
			}
			fmt.Fprintf(w, "- %s\n", desc)
		}
	}

	fmt.Fprintln(w, "\n---")
	fmt.Fprintln(w, "Use `ctx read <url>` to fetch full documents.")
}

func docSourceHeading(src docSource) string {
	title := strings.TrimSpace(src.title)
	if title == "" || strings.EqualFold(title, "unknown") {
		return src.url
	}
	return fmt.Sprintf("[%s](%s)", title, src.url)
}
