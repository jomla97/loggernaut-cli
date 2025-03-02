package cmd

import (
	"errors"
	"path"
	"slices"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Run: func(cmd *cobra.Command, args []string) {
		s := viper.GetStringSlice("sources")
		if len(s) == 0 {
			println("No sources configured")
			return
		}
		for _, source := range s {
			println(source)
		}
	},
}

// sourcesAddCmd represents the add subcommand of the sources command
var sourcesAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new source. Must be an absolute path to either a file or a folder.",
	Long:  `Add a new source. Must be an absolute path to either a file or a folder.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sources := viper.GetStringSlice("sources")
		for _, source := range args {
			if !path.IsAbs(source) {
				return errors.New("source must be an absolute path")
			}
			if !slices.Contains(sources, path.Dir(source)) && !slices.Contains(sources, source) {
				sources = append(sources, source)
			}
		}
		viper.Set("sources", sources)
		return viper.WriteConfig()
	},
}

// sourcesRemoveCmd represents the remove subcommand of the sources command
var sourcesRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a configured source.",
	Long:  `Remove a configured source.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configured := viper.GetStringSlice("sources")
		sources := []string{}
		for _, source := range configured {
			if !slices.Contains(args, source) {
				sources = append(sources, source)
			}
		}
		viper.Set("sources", sources)
		return viper.WriteConfig()
	},
}

// sourcesClearCmd represents the clear subcommand of the sources command
var sourcesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all configured sources.",
	Long:  `Remove all configured sources.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("sources", []string{})
		return viper.WriteConfig()
	},
}
