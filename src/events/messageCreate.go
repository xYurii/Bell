package events

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/handler"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" || m.Author.Bot {
		return
	}

	ctx := context.Background()
	guildData := database.Guild.GetGuild(ctx, m.GuildID)
	prefix := guildData.Prefix
	fmt.Printf("%+v", guildData)

	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	content := strings.TrimPrefix(m.Content, prefix)
	args := strings.Fields(strings.Trim(content, " "))

	if len(args) == 0 {
		return
	}

	commandName := args[0]
	command, exists := handler.GetCommand(commandName)

	if !exists {
		return
	}
	fmt.Printf("%s used %s command\n", m.Author.Username, commandName)

	command.Run(ctx, s, m, args[1:])
}
