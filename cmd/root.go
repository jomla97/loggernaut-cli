package cmd

import (
	"fmt"
	"os"

	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "loggernaut",
	Short: "Loggernaut is a CLI tool for collecting logs from various sources and sending them to the Loggernaut API.",
	// 	Long: `A longer description that spans multiple lines and likely contains
	// examples and usage of using your application. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Add subcommands to root command
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(sourcesCmd)
	rootCmd.AddCommand(collectCmd)

	// Add subcommands to sources command
	sourcesCmd.AddCommand(sourcesListCmd)
	sourcesCmd.AddCommand(sourcesAddCmd)
	sourcesCmd.AddCommand(sourcesRemoveCmd)
	sourcesCmd.AddCommand(sourcesClearCmd)

	// Execute root command
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	for _, dir := range []string{config.BasePath, config.OutboxPath} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Errorf("failed to make directory '%s': %w", dir, err))
		}
	}
	i, err := os.Stat(config.BasePath)
	if err != nil {
		panic(fmt.Errorf("failed to stat '%s': %w", config.BasePath, err))
	}
	if !i.IsDir() {
		panic(fmt.Errorf("'%s' is not a directory", config.BasePath))
	}
	viper.AddConfigPath(config.BasePath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetDefault("sources", []string{})
	viper.SafeWriteConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config: %w", err))
	}

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.loggernaut-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	sourcesAddCmd.Flags().StringSliceP("tags", "t", []string{}, "Tags to associate with the source")
	sourcesAddCmd.Flags().Bool("no-recursive", false, "Walk the source directory recursively")
}
