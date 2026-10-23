package shared

import "fmt"

type Config struct {
	ReportPath  string
	Includes    []string
	Excludes    []string
	ContextSize int
	Base        string
}

func NewConfig() *Config {
	return &Config{
		Base:        ".",
		Includes:    []string{},
		Excludes:    []string{},
		ContextSize: 3,
	}
}

func ValidateConfig(c *Config) error {
	if c.ContextSize < 1 {
		return fmt.Errorf("context-size must be a positive integer, got %d", c.ContextSize)
	}
	return nil
}
