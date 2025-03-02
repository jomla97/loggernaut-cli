package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/cobra"
)

// sourcesCmd represents the sources command
var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "Manage configured sources.",
	Long:  `Manage configured sources.`,
}

// sourcesListCmd represents the list subcommand of the sources command
var sourcesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured sources.",
	Long:  `List all configured sources.",`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the configured sources
		sources, err := config.GetAllSources()
		if err != nil {
			return err
		}

		// Check if there are no sources
		if len(sources) == 0 {
			println("No sources configured")
			return nil
		}

		// Print the sources
		for i, source := range sources {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("Source %d\n", i+1)
			fmt.Printf("System: %s\n", source.System)
			fmt.Printf("Path: %s\n", source.Path)
			fmt.Printf("Tags: %s\n", source.Tags)
		}
		return nil
	},
}

// sourcesAddCmd represents the add subcommand of the sources command
var sourcesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new source.",
	Long:  `Add a new source.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := config.Source{System: args[0], Path: args[1]}

		// Check if the source is valid
		if source.System == "" || source.Path == "" {
			return errors.New("system and path must be provided")
		} else if !path.IsAbs(os.ExpandEnv(source.Path)) {
			return errors.New("path must be an absolute path")
		}

		// Get the source tags
		if tags, err := cmd.Flags().GetStringSlice("tags"); err != nil {
			return err
		} else {
			source.Tags = tags
		}

		// Get the configured sources
		sources, err := config.GetAllSources()
		if err != nil {
			return err
		}

		// Check if the source path is already configured
		for _, s := range sources {
			if s.Path == source.Path || s.Path == path.Dir(source.Path) {
				return errors.New("source path already configured")
			}
		}

		// Write the updated sources to the config file
		return config.SetSources(append(sources, source))
	},
}

// sourcesRemoveCmd represents the remove subcommand of the sources command
var sourcesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a configured source.",
	Long:  `Remove a configured source.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the source path
		p, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}

		// Get the source index
		i, err := cmd.Flags().GetInt("index")
		if err != nil {
			return err
		}

		// Check that either the path or index flag is provided
		if p == "" && i <= 0 {
			return errors.New("either the path or index flag must be provided")
		}

		// Get the configured sources
		sources, err := config.GetAllSources()
		if err != nil {
			return err
		}

		// Remove the source
		var newSources []config.Source
		for j, source := range sources {
			if (p != "" && source.Path == p) || (i > 0 && i-1 == j) {
				continue
			}
			newSources = append(newSources, source)
		}

		if len(newSources) == len(sources) {
			return errors.New("source not found")
		}

		// Write the updated sources to the config file
		return config.SetSources(newSources)
	},
}

// sourcesClearCmd represents the clear subcommand of the sources command
var sourcesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all configured sources.",
	Long:  `Remove all configured sources.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Write an empty list of sources to the config file
		return config.SetSources([]config.Source{})
	},
}
