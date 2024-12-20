package commands

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"unsafe"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/events"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "botinfo",
		Aliases:     []string{"bot", "binfo"},
		Description: "BotInfo.Help",
		Cooldown:    3,
		Category:    "Info",
		Run:         runBotInfo,
	})
}

func runBotInfo(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		guild, err = s.Guild(m.GuildID)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	user := database.User.GetUser(ctx, m.Author)

	guildsCount := len(s.State.Guilds)

	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	ramUsage := float64(memory.Sys) / (1024 * 1024)

	guildUsage := 0
	if err == nil {
		guildUsage = int(unsafe.Sizeof(*guild))
		for _, member := range guild.Members {
			guildUsage += int(unsafe.Sizeof(*(member.User)))
			guildUsage += int(unsafe.Sizeof(*member))
		}
		for _, role := range guild.Roles {
			guildUsage += int(unsafe.Sizeof(*role))
		}
		for _, channel := range guild.Channels {
			guildUsage += int(unsafe.Sizeof(*channel))
		}
		for _, emoji := range guild.Emojis {
			guildUsage += int(unsafe.Sizeof(*emoji))
		}
		for _, presence := range guild.Presences {
			guildUsage += int(unsafe.Sizeof(*presence.Activities[0]))
			guildUsage += int(unsafe.Sizeof(*presence))
		}
		for _, voiceState := range guild.VoiceStates {
			guildUsage += int(unsafe.Sizeof(*voiceState))
		}
	}

	guildUsageText := services.Translate("BotInfo.GuildUsageKB", &user, map[string]interface{}{
		"GuildUsage": fmt.Sprintf("%.2f", float32(guildUsage)/1000),
	})
	if guildUsage > 1000000 {
		guildUsageText = services.Translate("BotInfo.GuildUsageMB", &user, map[string]interface{}{
			"GuildUsage": fmt.Sprintf("%.2f", float32(guildUsage)/1000/1000),
		})
	}

	ping := s.HeartbeatLatency().Milliseconds()

	myself := s.State.User

	avatar := discordgo.EndpointUserAvatar(myself.ID, myself.Avatar)
	myselfID, _ := strconv.ParseUint(myself.ID, 10, 64)
	creationDate := ((myselfID >> 22) + 1420070400000) / 1000
	readyAtMs := handler.ReadyAt.Unix()
	readyAt := fmt.Sprintf("<t:%d:F> (<t:%d:R>)", readyAtMs, readyAtMs)
	createdAt := fmt.Sprintf("<t:%d:F> (<t:%d:R>)", creationDate, creationDate)

	owner, err := s.User(events.OwnersIDs[1])
	ownerUsername := ""
	ownerID := ""
	if err == nil {
		ownerUsername = owner.Username
		ownerID = owner.ID
	}

	embed := discord.NewEmbedBuilder().
		WithTitle(services.Translate("BotInfo.Title", &user)).
		WithColor(utils.ColorDefault).
		WithFooter(fmt.Sprintf("@%s", m.Author.Username), m.Author.AvatarURL("2064")).
		WithField(
			services.Translate("BotInfo.Fields.Servers.Name", &user),
			services.Translate("BotInfo.Fields.Servers.Value", &user, map[string]interface{}{
				"GuildsCount": guildsCount,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.Ping.Name", &user),
			services.Translate("BotInfo.Fields.Ping.Value", &user, map[string]interface{}{
				"Ping": ping,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.RamUsage.Name", &user),
			services.Translate("BotInfo.Fields.RamUsage.Value", &user, map[string]interface{}{
				"RamUsage":   fmt.Sprintf("%.2f", ramUsage),
				"GuildUsage": guildUsageText,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.CreatedAt.Name", &user),
			services.Translate("BotInfo.Fields.CreatedAt.Value", &user, map[string]interface{}{
				"CreationDate": createdAt,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.ReadyAt.Name", &user),
			services.Translate("BotInfo.Fields.ReadyAt.Value", &user, map[string]interface{}{
				"ReadyAt": readyAt,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.Shard.Name", &user),
			services.Translate("BotInfo.Fields.Shard.Value", &user, map[string]interface{}{
				"Shard":  s.ShardID,
				"Shards": s.ShardCount,
			}), false).
		WithField(
			services.Translate("BotInfo.Fields.Workers.Name", &user),
			services.Translate("BotInfo.Fields.Workers.Value", &user, map[string]interface{}{
				"FreeInteractionsWorkers":  events.GetFreeWorkers(),
				"TotalInteractionsWorkers": events.Workers,
				"FreeMessagesWorkers":      events.GetFreeMessageWorkers(),
				"TotalMessagesWorkers":     events.MessageWorkers,
			}), false,
		).
		WithField(
			services.Translate("BotInfo.Fields.Owner.Name", &user),
			services.Translate("BotInfo.Fields.Owner.Value", &user, map[string]interface{}{
				"Owner": ownerUsername,
				"ID":    ownerID,
			}), false,
		).
		WithThumbnail(avatar)

	discord.NewMessage(s, m.ChannelID, m.ID).
		WithEmbed(embed.Build()).
		Send()

}
