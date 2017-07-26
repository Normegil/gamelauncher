package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

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
	flag.StringVar(&gamesFile, "data", "", "path to config file containing the list of video games (default is $HOME/.gamelauncher.yaml)")
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
		if "" != game.Command() && !game.Disabled() {
			commands = append(commands, &cobra.Command{
				Use: game.Command(),
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
		tree, ok := gamesTree.Get(game).(*toml.Tree)
		if !ok {
			return nil, errors.New("Game should be an instance of toml.Tree, got: " + reflect.TypeOf(tree).String())
		}

		name, ok := tree.Get("name").(string)
		if !ok {
			return nil, errors.New("'name' should be an instance of 'string', got: " + reflect.TypeOf(name).String())
		}

		command, ok := tree.Get("command").(string)
		if !ok && "" != command {
			return nil, errors.New("'command' should be an instance of 'string', got: " + reflect.TypeOf(command).String())
		}

		script, ok := tree.Get("script").(string)
		if !ok && "" != script {
			return nil, errors.New("'script' should be an instance of 'string', got: " + reflect.TypeOf(script).String())
		}

		scriptArgs := make([]string, 0)
		tmp := tree.Get("script-args")
		if nil != tmp {
			tmpArgs, ok := tmp.([]interface{})
			if !ok {
				return nil, errors.New("'script-args' should be an instance of '[]string', got: " + reflect.TypeOf(tmp).String())
			}
			for _, arg := range tmpArgs {
				strArg, ok := arg.(string)
				if !ok {
					return nil, errors.New("'script-args' should be an instances of 'string', got: " + reflect.TypeOf(arg).String())
				}
				scriptArgs = append(scriptArgs, strArg)
			}
		}

		disabled, ok := tree.Get("disabled").(bool)
		if !ok && false != disabled {
			return nil, errors.New("'disabled' should be an instance of 'bool', got: " + reflect.TypeOf(disabled).String())
		}

		games = append(games, model.NewGame(name, command, disabled, script, scriptArgs))
	}
	return games, nil
}
