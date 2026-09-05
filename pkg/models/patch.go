package models

// Represent a continuous sequence of source code
type Patch struct {
	File         string `json:"file"`
	Start        int    `json:"start"`
	End          int    `json:"end"`
	Source       string `json:"source"`
	SourceBefore string `json:"source_before"`
	SourceAfter  string `json:"source_after"`
}
