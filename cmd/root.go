package cmd

import (
	"flag"
	"fmt"
	"os"
	"sort"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/normegil/gamelauncher/model"
	toml "github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

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
	var gamesFile string
	flag.StringVar(&gamesFile, "data", "", "path to config file containing the list of video games (default is $HOME/.game.yaml)")
	flag.Parse()

	if "" == gamesFile {
		home, err := homedir.Dir()
		if err != nil {
			panic(err)
		}
		gamesFile = home + "/.games.toml"
	}

	games, err := loadGames(gamesFile)
	if err != nil {
		panic(err)
	}
	commands := make([]*cobra.Command, 0)
	for _, game := range games {
		if "" != game.Command && !game.Disabled {
			commands = append(commands, &cobra.Command{
				Use: game.Command,
				Run: getLaunchFunc(game),
			})
		}
	}
	RootCmd.AddCommand(commands...)
}

func getLaunchFunc(g *model.Game) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := g.Launch(); err != nil {
			panic(err)
		}
	}
}

func loadGames(gamesFile string) ([]*model.Game, error) {
	gamesTree, err := toml.LoadFile(gamesFile)
	if err != nil {
		panic(err)
	}

	for _, game := range gamesTree.Keys() {
		tree := gamesTree.Get(game).(*toml.Tree)

		var game model.Game
		err := tree.Unmarshal(&game)
		if err != nil {
			return nil, err
		}

		games = append(games, &game)
	}
	sort.Sort(model.ByName(games))
	return games, nil
}
