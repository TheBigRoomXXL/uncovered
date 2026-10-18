package shared

type Config struct {
	ReportPath  string
	Includes    []string
	Excludes    []string
	ContextSize int
	Base        string
}

// Coverage of every file in the reports, keyed by the normalized absolute path
type Coverage = map[string]FileCoverage

// FileCoverage holds the coverage data of a single file
type FileCoverage struct {
	File  string       // path as written in the report, used for display
	Lines map[int]bool // key is line number, value is covered or not
}
