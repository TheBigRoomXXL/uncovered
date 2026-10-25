package shared

import (
	"errors"
	"fmt"
	"io"
)

type Config struct {
	ReportPath  string
	Includes    []string
	Excludes    []string
	ContextSize int
	Base        string
	Ouput       io.Writer
	Render      Renderer
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

	if c.Ouput == nil {
		return errors.New("ouput not set")
	}

	if c.Render == nil {
		return errors.New("render function not set")
	}

	return nil
}
