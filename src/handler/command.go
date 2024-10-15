package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/prototypes"
)

type Command struct {
	Name        string
	Description string
	Aliases     []string
	Usage       string
	Category    string
	Developer   bool
	Cooldown    int
	Run         func(context.Context, *discordgo.Session, *discordgo.MessageCreate, []string)
}

var Commands = map[string]Command{}

func RegisterCommand(cmd Command) {
	Commands[cmd.Name] = cmd
}

func GetCommand(name string) (Command, bool) {
	for _, cmd := range Commands {
		if cmd.Name == name || prototypes.Includes(cmd.Aliases, func(alias string) bool {
			return alias == name
		}) {
			return cmd, true
		}
	}
	return Command{}, false
}
