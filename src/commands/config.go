package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:     "config",
		Aliases:  []string{"configurar", "configurações", "configuracao"},
		Run:      runConfig,
		Category: "Config",
	})
}

func runConfig(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	memberPermissions, _ := discord.MemberHasPermission(s, m.Message, discordgo.PermissionManageServer)

	user := database.User.GetUser(ctx, m.Author)
	guild := database.Guild.GetGuild(ctx, m.GuildID)
	reply := discord.NewMessage(s, m.ChannelID, m.ID)

	if !memberPermissions {
		response := services.Translate("Config.NoPermission", &user, m.Author.Mention())
		reply.WithContent(response).Send()
		return
	}

	var channels string
	if len(guild.CommandsChannels) > 0 {
		for _, id := range guild.CommandsChannels {
			channels += fmt.Sprintf("<#%s> (%s)\n", id, id)
		}
	} else {
		channels = services.Translate("Config.NoChannels", &user)
	}

	discordGuild, _ := s.State.Guild(m.GuildID)

	embed := createEmbedConfig(discordGuild, user, channels)

	reply.WithChannelSelect("config-channels", "Selecionar canais", 1, 15, false)

	msg, _ := reply.WithEmbed(embed).Send()

	handler.CreateMessageComponentCollector(msg, func(i *discordgo.Interaction) {
		handler.DeferInteraction(s, i)

		if i.Member.User.ID != m.Author.ID {
			return
		}

		chs := i.MessageComponentData().Values
		var c string
		if len(chs) > 0 {
			for _, id := range chs {
				c += fmt.Sprintf("<#%s> (%s)\n", id, id)
			}
		} else {
			c = services.Translate("Config.NoChannels", &user)
		}

		database.Guild.UpdateGuild(ctx, i.GuildID, func(g schemas.Guild) schemas.Guild {
			g.CommandsChannels = chs
			return g
		})

		newEmbed := createEmbedConfig(discordGuild, user, c)
		msgEdit := &discordgo.MessageEdit{
			Embeds:  &[]*discordgo.MessageEmbed{newEmbed},
			ID:      i.Message.ID,
			Channel: i.ChannelID,
		}

		_, err := s.ChannelMessageEditComplex(msgEdit)

		if err != nil {
			log.Fatalln("Error while editing message:", err)
		}

	}, 3*time.Minute)

}

func createEmbedConfig(guild *discordgo.Guild, user schemas.User, channels string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Footer: &discordgo.MessageEmbedFooter{
			Text:    guild.Name,
			IconURL: guild.IconURL("2048"),
		},
		Title: services.Translate("Config.EmbedTitle", &user),
		Color: utils.ColorDefault,
		Description: services.Translate("Config.Description", &user, map[string]interface{}{
			"Channels": channels,
		}),
	}
}
