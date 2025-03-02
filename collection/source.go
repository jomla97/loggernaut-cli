package collection

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/viper"
)

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

// Source represents a log source
type Source struct {
	System    string   `json:"system"`
	Path      string   `json:"path"`
	Tags      []string `json:"tags"`
	Recursive bool     `json:"recursive"`
}

// Walk walks the source path and returns the path to all log files within it
func (s *Source) Walk() ([]string, error) {
	var files []string
	err := filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk '%s': %w", path, err)
		}

		if info.IsDir() && !s.Recursive {
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".log" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return []string{}, err
	}
	return files, nil
}

func (s *Source) Collect() (int, error) {
	files, err := s.Walk()
	if err != nil {
		return 0, err
	}

	var moved int
	for _, path := range files {
		err := s.moveToOutbox(path)
		if err != nil {
			return moved, err
		}
		moved++
	}

	return moved, nil
}

// moveToOutbox moves the log file with the specified path to the outbox folder.
func (s *Source) moveToOutbox(srcPath string) error {
	// Open the source file
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}

	// Create the destination file
	dstPath := filepath.Join(config.OutboxPath, uuid.NewString()+filepath.Ext(srcPath))
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %w", err)
	}
	defer dst.Close()

	// Copy the file to the outbox
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file to outbox: %w", err)
	}
	src.Close()

	// Create the meta data file
	meta := Meta{Source: *s, OriginalPath: srcPath, OutboxPath: dstPath}
	if err := meta.Create(); err != nil {
		return err
	}

	// Remove the source file
	err = os.Remove(srcPath)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %w", err)
	}

	fmt.Printf("Moved %s to outbox\n", srcPath)
	return nil
}
