package events

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
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

	userData := database.User.GetUser(ctx, m.Author)
	guildData := database.Guild.GetGuild(ctx, m.GuildID)
	prefix := strings.ToLower(guildData.Prefix)

	if mentionsClient(s, m) {
		res := services.Translate("Bot.Mention", &userData, prefix)
		_, err := discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	if !strings.HasPrefix(strings.ToLower(m.Content), strings.ToLower(prefix)) {
		return
	}

	content := strings.TrimPrefix(m.Content, prefix)
	args := strings.Fields(strings.Trim(content, " "))

	if len(args) == 0 {
		return
	}

	commandName := strings.ToLower(args[0])
	command, exists := handler.GetCommand(commandName)

	if !exists {
		return
	}

	if command.Developer && !prototypes.Includes(OwnersIDs, func(id string) bool {
		return m.Author.ID == id
	}) {
		return
	}

	if len(guildData.CommandsChannels) > 0 {
		inAllowedChannel := prototypes.Includes(guildData.CommandsChannels, func(id string) bool {
			return id == m.ChannelID
		})
		if !inAllowedChannel {
			hasChannel := false
			allowedExistentChannel := services.Translate("Config.NoChannelExistent", &userData, map[string]interface{}{
				"Prefix": guildData.Prefix,
			})
			for _, id := range guildData.CommandsChannels {
				channel, err := s.State.Channel(id)
				if err == nil {
					allowedExistentChannel = channel.Mention()
					hasChannel = true
					break
				}
			}

			if !hasChannel {
				database.Guild.UpdateGuild(ctx, m.GuildID, func(g schemas.Guild) schemas.Guild {
					g.CommandsChannels = []string{}
					return g
				})
				discord.NewMessage(s, m.ChannelID, m.ID).WithContent(services.Translate("Config.TryAgain", &userData)).Send()
				return
			}

			res := services.Translate("Config.WarnMessage", &userData, map[string]interface{}{
				"AllowedChannel": allowedExistentChannel,
			})
			warnMsg, _ := discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
			time.Sleep(10 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, warnMsg.ID)
			return
		}
	}

	database.Commands.InsertCommand(ctx, commandName, m.Author.ID, m.GuildID, m.ChannelID, m.Content)

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

func mentionsClient(s *discordgo.Session, evt *discordgo.MessageCreate) bool {
	if evt.ReferencedMessage != nil {
		return false
	}

	if len(evt.Mentions) == 1 {
		return evt.Mentions[0].ID == s.State.User.ID
	}
	return false
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
