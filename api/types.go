package api

type librarySearchResponse struct {
	Results []Library `json:"results"`
	Error   string    `json:"error,omitempty"`
}

type Library struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	TotalSnippets  int      `json:"totalSnippets"`
	TotalTokens    int      `json:"totalTokens,omitempty"`
	Stars          int      `json:"stars,omitempty"`
	TrustScore     float64  `json:"trustScore,omitempty"`
	BenchmarkScore float64  `json:"benchmarkScore,omitempty"`
	Versions       []string `json:"versions,omitempty"`
}

type DocsResponse struct {
	CodeSnippets []CodeSnippet `json:"codeSnippets"`
	InfoSnippets []InfoSnippet `json:"infoSnippets"`
	Error        string        `json:"error,omitempty"`
}

type CodeSnippet struct {
	CodeTitle       string `json:"codeTitle"`
	CodeDescription string `json:"codeDescription"`
	CodeID          string `json:"codeId"`
	PageTitle       string `json:"pageTitle"`
}

type InfoSnippet struct {
	URL        string `json:"url,omitempty"`
	Breadcrumb string `json:"breadcrumb,omitempty"`
	Content    string `json:"content"`
}
