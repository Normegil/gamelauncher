package cmd

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/normegil/gamelauncher/model"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

var gamesFile string
var games []*model.Game

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
	if "" == gamesFile {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		gamesFile = home + "/.games.toml"
	}
	gamesTree, err := toml.LoadFile(gamesFile)
	if err != nil {
		panic(err)
	}

	for _, game := range gamesTree.Keys() {
		tree, ok := gamesTree.Get(game).(*toml.Tree)
		if !ok {
			panic(errors.New("Game should be an instance of toml.Tree, got: " + reflect.TypeOf(tree).String()))
		}

		name, ok := tree.Get("name").(string)
		if !ok {
			panic(errors.New("'name' should be an instance of 'string', got: " + reflect.TypeOf(name).String()))
		}

		command, ok := tree.Get("command").(string)
		if !ok {
			panic(errors.New("'command' should be an instance of 'string', got: " + reflect.TypeOf(command).String()))
		}

		games = append(games, model.NewGame(name, command))
	}
}
