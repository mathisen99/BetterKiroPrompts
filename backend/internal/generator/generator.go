package generator

// GeneratedFile represents a single generated output file.
type GeneratedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}
