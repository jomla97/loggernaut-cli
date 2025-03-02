package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect logs from the configured sources",
	Long:  `Collect logs from the configured sources`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sources, err := config.GetAllSources()
		if err != nil {
			return err
		}

		collected, err := collect(sources)
		if err != nil {
			return err
		}

		if collected == 0 {
			fmt.Println("Found no log files to collect.")
			return nil
		}

		fmt.Printf("Successfully collected %d log files.\n", collected)
		return nil
	},
}

// collect logs from the configured sources, moving them to the outbox folder. Returns the number of files collected.
func collect(sources []config.Source) (int, error) {
	var collected int
	for _, source := range sources {
		paths, err := walkSourcePath(source)
		if err != nil {
			return collected, err
		}

		moved, err := moveFiles(source, paths)
		if err != nil {
			return collected, err
		}
		collected += moved
	}
	return collected, nil
}

func walkSourcePath(source config.Source) ([]string, error) {
	var files []string
	err := filepath.Walk(source.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk '%s': %w", path, err)
		}

		if info.IsDir() && !source.Recursive {
			return filepath.SkipDir
		}

		if filepath.Ext(path) == ".log" {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// moveFiles moves the log files to the outbox folder, returning the number of files moved.
func moveFiles(source config.Source, paths []string) (int, error) {
	var moved int
	for _, path := range paths {
		err := moveFile(source, path)
		if err != nil {
			return moved, err
		}
		moved++
	}
	return moved, nil
}

// moveFile moves the log file to the outbox folder.
func moveFile(source config.Source, srcPath string) error {
	// Open the source file
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}

	// Create the destination file
	dstPath := filepath.Join(config.OutboxPath, uuid.NewString()+filepath.Ext(srcPath))
	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to open destination file: %v", err)
	}
	defer dst.Close()

	// Copy the file to the outbox
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file to outbox: %v", err)
	}

	src.Close()

	// Remove the source file
	err = os.Remove(srcPath)
	if err != nil {
		return fmt.Errorf("failed to remove source file: %v", err)
	}

	if err := createMetaFile(source, srcPath, dstPath); err != nil {
		return err
	}

	fmt.Printf("Moved %s to outbox\n", srcPath)
	return nil
}

func createMetaFile(source config.Source, srcPath, dstPath string) error {
	// Create the meta file
	filename := fmt.Sprintf("%s.meta.json", filepath.Base(dstPath))
	file, err := os.Create(filepath.Join(config.OutboxPath, filename))
	if err != nil {
		return fmt.Errorf("failed to create meta file: %v", err)
	}
	defer file.Close()

	meta := struct {
		Source       config.Source `json:"source"`
		OriginalPath string        `json:"path"`
	}{
		Source:       source,
		OriginalPath: srcPath,
	}

	// Marshal the meta data to json
	b, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("failed to marshal meta data: %v", err)
	}

	// Write the meta data to the file
	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write to meta file: %v", err)
	}

	return nil
}
