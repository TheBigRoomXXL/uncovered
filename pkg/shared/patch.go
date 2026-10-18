package shared

// Represent a continuous sequence of source code
type Patch struct {
	File         string   `json:"file"`
	FileAbs      string   `json:"file_abs"`
	Start        int      `json:"start"` // Inclusive
	End          int      `json:"end"`   // Inclusive
	Source       []string `json:"source"`
	SourceBefore []string `json:"source_before"`
	SourceAfter  []string `json:"source_after"`
	SourceError  error    `json:"source_error"`
}
