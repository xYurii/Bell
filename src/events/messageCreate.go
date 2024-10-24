package events

import (
	"context"
	"log"
	"strings"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/prototypes"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils/discord"
)

const MessageWorkers = 128

var MessageCreateChannel = make(chan *discordgo.MessageCreate, MessageWorkers)
var freeMessageWorkersCounter int64 = MessageWorkers

func HandleMessageCreate(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.GuildID == "" || m.Author.Bot {
		return
	}

	guildData := database.Guild.GetGuild(ctx, m.GuildID)
	prefix := strings.ToLower(guildData.Prefix)

	mentions := m.Mentions
	if len(mentions) == 1 && mentions[0].ID == s.State.User.ID {
		userData := database.User.GetUser(ctx, m.Author)
		res := services.Translate("Bot.Mention", &userData, prefix)
		_, err := discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
		if err != nil {
			log.Fatalln(err)
		}
		return
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

	if command.Developer && !prototypes.Includes(OwnersIDs, func(id string) bool {
		return m.Author.ID == id
	}) {
		return
	}

	command.Run(ctx, s, m, args[1:])
}

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ctx := context.Background()
	// go HandleMessageCreate(ctx, s, m)
	select {
	case MessageCreateChannel <- m:
	default:
		log.Printf("MessageCreateChannel is full, dropping message %s", m.ID)
	}
}

func MessageWorker(s *discordgo.Session) {
	for msg := range MessageCreateChannel {
		atomic.AddInt64(&freeMessageWorkersCounter, -1)
		ctx := context.Background()
		HandleMessageCreate(ctx, s, msg)
		atomic.AddInt64(&freeMessageWorkersCounter, 1)
	}
}

func InitMessageWorkers(s *discordgo.Session) {
	for i := 0; i < MessageWorkers; i++ {
		go MessageWorker(s)
	}
}

func GetFreeMessageWorkers() int64 {
	return atomic.LoadInt64(&freeMessageWorkersCounter)
}

// func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	if m.GuildID == "" || m.Author.Bot {
// 		return
// 	}

// 	ctx := context.Background()
// 	guildData := database.Guild.GetGuild(ctx, m.GuildID)
// 	prefix := strings.ToLower(guildData.Prefix)

// 	mentions := m.Mentions
// 	if len(mentions) == 1 && mentions[0].ID == s.State.User.ID {
// 		userData := database.User.GetUser(ctx, m.Author)
// 		res := services.Translate("Bot.Mention", &userData, prefix)
// 		_, err := discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
// 		if err != nil {
// 			log.Fatalln(err)
// 		}
// 	}

// 	if !strings.HasPrefix(strings.ToLower(m.Content), prefix) {
// 		return
// 	}

// 	content := strings.TrimPrefix(m.Content, prefix)
// 	args := strings.Fields(strings.Trim(content, " "))

// 	if len(args) == 0 {
// 		return
// 	}

// 	commandName := args[0]
// 	command, exists := handler.GetCommand(commandName)

// 	if !exists {
// 		return
// 	}

// 	if command.Developer && !prototypes.Includes(OwnersIDs, func(id string) bool {
// 		return m.Author.ID == id
// 	}) {
// 		return
// 	}

// 	command.Run(ctx, s, m, args[1:])
// }
