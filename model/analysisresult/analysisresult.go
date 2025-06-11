package analysisresult

type AnalysisResult struct {
	Title             string         `json:"title"`
	Headings          map[string]int `json:"headings"`
	HTMLVersion       string         `json:"html_version"`
	InternalLinks     []string       `json:"internal_links"`
	ExternalLinks     []string       `json:"external_links"`
	InaccessibleLinks []string       `json:"inaccessible_links"`
	LoginForm         bool           `json:"is_login_form"`
}
