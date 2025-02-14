package components

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterComponent(handler.Component{
		Name: "guild",
		Type: discordgo.ButtonComponent,
		Run:  notifyGuild,
	})
}

func notifyGuild(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	handler.DeferInteraction(s, i.Interaction)

	channelId := os.Getenv("CHANNEL_NOTIFY")

	component := i.MessageComponentData()
	choice := strings.ReplaceAll(component.CustomID[strings.IndexByte(component.CustomID, '_')+1:], "-", " ")

	embed := discord.NewEmbedBuilder().
		WithAuthor(i.User.Username, i.User.AvatarURL("2048"), "").
		WithDescription(fmt.Sprintf("O usuário conheceu o bot/servidor por: %s", choice)).
		WithColor(utils.ColorDefault)

	reply := discord.NewMessage(s, channelId, "").WithEmbed(embed.Build())

	reply.Send()

	msgEdit := &discordgo.MessageEdit{
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Components: &[]discordgo.MessageComponent{},
	}

	_, err := s.ChannelMessageEditComplex(msgEdit)

	if err != nil {
		log.Fatalln("Error while editing message:", err)
	}
}
