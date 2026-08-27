package models

type Config struct {
	reportPath  string
	includes    []string
	excludes    []string
	contextSize int
}
