package commands

import (
	"context"
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/xYurii/Bell/src/handler"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "eval",
		Aliases:     []string{"ev"},
		Cooldown:    5,
		Run:         runEval,
		Category:    "general",
		Description: "Ping.Description",
		Developer:   true,
	})
}

func runEval(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	code := ""
	for _, arg := range args {
		code += arg + " "
	}

	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	exports := interp.Exports{
		"env": map[string]reflect.Value{
			"Author":  reflect.ValueOf(m.Author),
			"Content": reflect.ValueOf(m.Content),
			"Channel": reflect.ValueOf(m.ChannelID),
			"Session": reflect.ValueOf(s),
			"Message": reflect.ValueOf(m),
		},
	}
	i.Use(exports)

	result, err := i.Eval(code)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Erro ao avaliar código: %s", err.Error()))
		return
	}

	response := fmt.Sprintf("```rust\n%v\n```", result)
	s.ChannelMessageSend(m.ChannelID, response)
}
