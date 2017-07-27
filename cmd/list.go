package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var tagsFilter string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered games",
	Long:  `list all registered games`,
	Aliases: []string{
		"ls",
	},
	Run: func(cmd *cobra.Command, args []string) {
		tags := make([]string, 0)
		if "" != tagsFilter {
			tags = strings.Split(tagsFilter, ",")
		}
		for _, game := range games {
			if !game.Disabled && containAll(game.Tags, tags) {
				fmt.Printf("%20s\t%15s\t\t\t%+v\n", game.Name, "("+game.Command+")", game.Tags)
			}
		}
	},
}

func containAll(tested []string, toContain []string) bool {
	for _, tag := range toContain {
		var contains bool
		for _, tested := range tested {
			if tested == tag {
				contains = true
				break
			}
		}
		if !contains {
			return false
		}
	}
	return true
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVarP(&tagsFilter, "tags", "t", "", "Tags filter")
}
