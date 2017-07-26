package model

func NewGame(name string, command string) *Game {
	return &Game{
		name:    name,
		command: command,
	}
}

type Game struct {
	name    string
	command string
}

func (g Game) Name() string {
	return g.name
}

func (g Game) Command() string {
	return g.command
}

func (g Game) Launch() error {
	return nil
}
