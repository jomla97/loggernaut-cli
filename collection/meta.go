package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jomla97/loggernaut-cli/config"
)

// ReadMetaFile reads the data from the meta file for the given outbox log file
func ReadMetaFile(path string) (Meta, error) {
	metaPath := path + ".meta.json"
	meta := Meta{OutboxPath: path, MetaPath: &metaPath}

	// Open the meta data file
	file, err := os.Open(metaPath)
	if err != nil {
		return meta, fmt.Errorf("failed to open meta data file: %w", err)
	}

	// Copy the file data to a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return meta, fmt.Errorf("failed to copy data from meta data file: %w", err)
	}

	// Unmarshal the meta data
	err = json.Unmarshal(buf.Bytes(), &meta)
	if err != nil {
		return meta, fmt.Errorf("failed to unmarshal meta data: %w", err)
	}
	return meta, nil
}

type Meta struct {
	Source       Source  `json:"source"`
	OriginalPath string  `json:"path"`
	OutboxPath   string  `json:"-"`
	MetaPath     *string `json:"-"`
}

// createMetaFile creates a meta file for the given log file
func (m *Meta) Create() error {
	// Create the meta file
	filename := fmt.Sprintf("%s.meta.json", filepath.Base(m.OutboxPath))
	path := filepath.Join(config.OutboxPath, filename)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create meta file: %w", err)
	}
	defer file.Close()

	// Marshal the meta data to json
	b, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal meta data: %w", err)
	}

	// Write the meta data to the file
	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to meta file: %w", err)
	}

	// Update the meta file path in the meta struct
	m.MetaPath = &path

	return nil
}
