package handler

import (
	"context"
	"time"

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

var UserStatusTracking = make(map[string]int64)
var Commands = map[string]Command{}
var ReadyAt time.Time

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
