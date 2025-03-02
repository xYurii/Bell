package components

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
)

func init() {
	handler.RegisterComponent(handler.Component{
		Name: "allRoostersNames",
		Type: discordgo.ButtonComponent,
		Run:  allRoostersNames,
	})
}

func allRoostersNames(_ context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	roosters, err := utils.GetRoostersClasses("https://raw.githubusercontent.com/Acnologla/asura-site/main/public/resources/class.json")

	if err != nil {
		res := "Não foi possível obter os nomes dos galos!"
		handler.RespondInteraction(s, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
		return
	}

	allRoosters := utils.GetRoostersNames(roosters)

	var description string
	color := 0
	for i, rooster := range allRoosters {
		if i != 0 && i > 0 {
			description += fmt.Sprintf("**%s** - **%s**\n", rooster.Name, rooster.Rarity.String())
			if rooster.Rarity >= 5 {
				color = rooster.Rarity.Color()
			}
		}
	}

	embeds := &discordgo.MessageEmbed{
		Title:       "**Lista de Galos**",
		Description: description,
		Color:       color,
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embeds},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}

	s.InteractionRespond(i.Interaction, response)
}
