package cmd

import (
	"fmt"

	"github.com/jomla97/loggernaut-cli/collection"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send collected logs to the Loggernaut API",
	Long:  `Send collected logs to the Loggernaut API`,
	RunE: func(cmd *cobra.Command, args []string) error {
		debug, _ := cmd.Flags().GetBool("debug")
		sent, err := collection.SendAll(debug)
		if err != nil {
			return err
		}
		if sent == 0 {
			fmt.Println("No logs to send")
		} else {
			fmt.Printf("Sent %d log file(s).\n", sent)
		}
		return nil
	},
}
