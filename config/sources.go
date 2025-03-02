package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Source represents a log source
type Source struct {
	System string
	Path   string
	Tags   []string
}

// GetAllSources returns all configured sources
func GetAllSources() ([]Source, error) {
	var sources []Source
	if err := viper.UnmarshalKey("sources", &sources); err != nil {
		return sources, fmt.Errorf("failed to get sources: %w", err)
	}
	return sources, nil
}

// SetSources sets the configured sources
func SetSources(sources []Source) error {
	viper.Set("sources", sources)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("failed to update sources: %w", err)
	}
	return nil
}
