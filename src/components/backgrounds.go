package components

import (
	"context"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/commands"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
)

func init() {
	handler.RegisterComponent(handler.Component{
		Name: "backgrounds",
		Type: discordgo.ButtonComponent,
		Run:  cosmectics,
	})
}

func cosmectics(_ context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := handler.DeferInteraction(s, i.Interaction)
	if err != nil {
		log.Println("Error while deferring interaction:", err)
		return
	}

	backgrounds := utils.GetBackgrounds()
	customIdSplitted := handler.ParseComponentId(i.MessageComponentData().CustomID)
	action := customIdSplitted[1]
	userID := customIdSplitted[3]
	newPage, err := strconv.Atoi(customIdSplitted[2])

	if err != nil {
		log.Fatal(err)
		return
	}

	if userID != i.Member.User.ID {
		return
	}

	switch action {
	case "next":
		newPage++
	case "previous":
		newPage--
	}

	if newPage < 0 {
		newPage = len(backgrounds) - 1
	} else if newPage >= len(backgrounds) {
		newPage = 0
	}

	background := backgrounds[newPage]
	newEmbed := commands.BuildEmbed(background, newPage, len(backgrounds))
	newButtons := commands.CreateButtons(newPage, userID)

	msgEdit := &discordgo.MessageEdit{
		Embeds:     &[]*discordgo.MessageEmbed{newEmbed},
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
		Components: &newButtons,
	}

	_, err = s.ChannelMessageEditComplex(msgEdit)

	if err != nil {
		log.Fatalln("Error while editing message:", err)
	}
}
