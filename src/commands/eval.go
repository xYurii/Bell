package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/acnologla/interpreter"
	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

func init() {
	interpreter.Init(map[string]interface{}{
		"commands": handler.Commands,
	})

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

	msg_json, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error serializing message: %v", err))
		return
	}

	eval := interpreter.Run(code, map[string]interface{}{
		"msg": string(msg_json),
		"s":   s,
	})
	response := fmt.Sprintf("```rust\n%v\n```", eval)
	s.ChannelMessageSend(m.ChannelID, response)
}
