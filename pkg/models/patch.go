package models

// Represent a continuous sequence of source code
type Patch struct {
	File  string
	Start int
	End   int
}
