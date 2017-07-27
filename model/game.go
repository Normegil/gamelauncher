package model

import (
	"fmt"
	"os"
	"os/exec"
)

type Game struct {
	Name     string
	Command  string
	Disabled bool
	Tags     []string

	Script     string
	ScriptArgs []string
}

func (g Game) Launch() error {
	fmt.Println("Launching: " + g.Name)
	if "" == g.Script {
		fmt.Printf("Cannot launch %s: script is empty", g.Name)
		return nil
	}
	cmd := exec.Command(g.Script, g.ScriptArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

type ByName []*Game

func (b ByName) Len() int { return len(b) }
func (b ByName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByName) Less(i, j int) bool {
	return b[i].Name < b[j].Name
}
