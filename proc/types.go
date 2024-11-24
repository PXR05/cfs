package proc

type InputData struct {
	Items []string `json:"items"`
}

type Category struct {
	Name      string              `json:"name"`
	Keywords  []string            `json:"keywords"`
	Phrases   []string            `json:"phrases"`
	Contexts  map[string][]string `json:"contexts"`
	Excluders []string            `json:"excluders"`
}

type CategoryOutputData struct {
	Categories []Category `json:"categories"`
}

type ClassificationResult struct {
	Item       string   `json:"item"`
	Category   string   `json:"category"`
	Confidence float64  `json:"confidence"`
	Matches    []string `json:"matches"`
}

type ClassificationOutputData struct {
	Results []ClassificationResult `json:"results"`
}
