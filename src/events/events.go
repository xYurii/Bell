package events

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

var mutex sync.RWMutex

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" || m.Author.Bot {
		return
	}

	ctx := context.Background()
	prefix := ".."

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
