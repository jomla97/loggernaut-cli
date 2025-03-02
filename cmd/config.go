package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//http://localhost:80/log

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration.",
	Long:  `Manage configuration.`,
}

// configSetCmd represents the set subcommand of the config command
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value. Key can be one of 'api-url'.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "api-url":
			viper.Set("api_url", strings.TrimSuffix(args[1], "/"))
		default:
			return errors.New("invalid key")
		}
		return viper.WriteConfig()
	},
}

// configGetCmd represents the get subcommand of the config command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value. Can be one of 'api-url'.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "api-url":
			println(viper.GetString("api_url"))
			return nil
		}
		return errors.New("invalid key")
	},
}
