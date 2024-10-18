package events

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils/discord"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" || m.Author.Bot {
		return
	}

	ctx := context.Background()
	guildData := database.Guild.GetGuild(ctx, m.GuildID)
	prefix := strings.ToLower(guildData.Prefix)

	mentions := m.Mentions
	if len(mentions) == 1 && mentions[0].ID == s.State.User.ID {
		_, err := discord.NewMessage(s, m.ChannelID, m.ID).WithContent(fmt.Sprintf("Olá %s! Meu prefixo é **%s**\nPara ver meus comandos, digite **%shelp**!", m.Author.Mention(), prefix, prefix)).Send()
		if err != nil {
			log.Fatalln(err)
		}
	}

	if !strings.HasPrefix(strings.ToLower(m.Content), prefix) {
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
