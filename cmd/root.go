package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var gamesFile string
var games []Game

var RootCmd = &cobra.Command{
	Use:   "gamelauncher",
	Short: "List and launch registered games",
	Long:  `List and launch registered games.`,
	Run: func(cmd *cobra.Command, args []string) {
		listCmd.Run(cmd, args)
	},
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&gamesFile, "data", "", "path to config file containing the list of video games (default is $HOME/.gamelauncher.yaml)")
}

func initConfig() {
}
