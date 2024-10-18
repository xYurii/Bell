package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
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
	if len(cmd.Aliases) > 0 {
		for _, alias := range cmd.Aliases {
			if alias != cmd.Name {
				Commands[alias] = cmd
			}
		}
	}
}

func GetCommand(name string) (Command, bool) {
	cmd, exists := Commands[name]
	return cmd, exists
}
