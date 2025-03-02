package cmd

import (
	"errors"

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
	Use:   "set",
	Short: "Set a configuration value.",
	Long:  `Set a configuration value.",`,
}

// configGetCmd represents the get subcommand of the config command
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value for the specified key. Can be one of 'api-url'.",
	Long:  `Get a configuration value for the specified key. Can be one of 'api-url'.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		switch key {
		case "api-url":
			println(viper.GetString("api_url"))
			return nil
		}
		return errors.New("invalid key")
	},
}

// configSetApiUrlCmd represents the set-api-url subcommand of the set subcommand
var configSetApiUrlCmd = &cobra.Command{
	Use:   "api-url <url>",
	Short: "Set the Loggernaut API URL.",
	Long:  `Set the Loggernaut API URL.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("api_url", args[0])
		return viper.WriteConfig()
	},
}
