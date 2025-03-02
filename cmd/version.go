package cmd

import (
	"fmt"

	"github.com/jomla97/loggernaut-cli/config"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Loggernaut CLI",
	Long:  `All software has versions. This is Loggernaut CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Loggernaut CLI %s\n", config.Version)
	},
}
