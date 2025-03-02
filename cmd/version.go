package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Loggernaut CLI",
	Long:  `All software has versions. This is Loggernaut CLI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Loggernaut CLI %s\n", version)
	},
}
