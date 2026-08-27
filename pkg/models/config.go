package models

type Config struct {
	ReportPath  string
	Includes    []string
	Excludes    []string
	ContextSize int
}
