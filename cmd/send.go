package cmd

import (
	"github.com/jomla97/loggernaut-cli/collection"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send collected logs to the Loggernaut API",
	Long:  `Send collected logs to the Loggernaut API`,
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		err := collection.SendAll(debug)
		if err != nil {
			return err
		}
		return nil
	},
}
