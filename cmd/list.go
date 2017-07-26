package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered games",
	Long:  `list all registered games`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, game := range games {
			fmt.Printf("%s\t\t\t(%s)\n", game.Name(), game.Command())
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
