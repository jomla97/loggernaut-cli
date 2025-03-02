package cmd

import (
	"fmt"

	"github.com/jomla97/loggernaut-cli/collection"
	"github.com/spf13/cobra"
)

var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "Collect logs from the configured sources",
	Long:  `Collect logs from the configured sources`,
	RunE: func(cmd *cobra.Command, args []string) error {
		sources, err := collection.GetAllSources()
		if err != nil {
			return err
		}

		collected, err := collection.CollectAll(sources)
		if err != nil {
			return err
		}

		if collected == 0 {
			fmt.Println("Found no log files to collect.")
			return nil
		}

		fmt.Printf("Successfully collected %d log file(s).\n", collected)
		return nil
	},
}
