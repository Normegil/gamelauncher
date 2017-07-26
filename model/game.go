package model

import (
	"fmt"
	"os"
	"os/exec"
)

func NewGame(name string, command string, disabled bool, script string, scriptArgs []string) *Game {
	var cmd *exec.Cmd
	if "" != script {
		cmd = exec.Command(script, scriptArgs...)
	}
	return &Game{
		name:     name,
		command:  command,
		cmd:      cmd,
		disabled: disabled,
	}
}

type Game struct {
	name     string
	command  string
	cmd      *exec.Cmd
	disabled bool
}

func (g Game) Name() string {
	return g.name
}

func (g Game) Command() string {
	return g.command
}

func (g Game) Disabled() bool {
	return g.disabled
}

func (g *Game) Launch() error {
	fmt.Println("Launching: " + g.Name())
	if nil == g.cmd {
		fmt.Printf("Cannot launch %s: script is empty", g.Name())
		return nil
	}
	g.cmd.Stdout = os.Stdout
	g.cmd.Stderr = os.Stderr
	g.cmd.Stdin = os.Stdin
	return g.cmd.Run()
}
