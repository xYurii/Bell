package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/prototypes"
	"github.com/xYurii/Bell/src/utils"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:     "galo",
		Aliases:  []string{"rooster"},
		Cooldown: 5,
		Run:      runAsuraRooster,
	})
}

func runAsuraRooster(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	response := &discordgo.MessageSend{
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	}

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Skins",
					Style:    discordgo.PrimaryButton,
					CustomID: "skins",
				},
				discordgo.Button{
					Label:    "Habilidades",
					Style:    discordgo.SecondaryButton,
					CustomID: "skills",
				},
			},
		},
	}

	if len(args) == 0 {
		response.Content = "Você não informou o nome do galo que deseja visualizar!"
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	roosterName := args[0]

	if roosterName == "" {
		response.Content = "Você deve informar o nome do galo que deseja visualizar!"
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	roosters, err := utils.GetJSON("https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/class.json")

	if err != nil {
		response.Content = "Eu não consegui buscar a lista dos nomes dos galos... Tente novamente."
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	allValidRoosters := utils.GetRoostersNames(roosters)
	roosterExists := prototypes.Includes(allValidRoosters, func(roosterData utils.Class) bool {
		return strings.EqualFold(roosterData.Name, roosterName)
	})

	if !roosterExists {
		response.Content = fmt.Sprintf("O galo **%s** não existe!", roosterName)
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	rooster := prototypes.Find(allValidRoosters, func(rooster utils.Class) bool {
		return strings.EqualFold(rooster.Name, roosterName)
	})

	roostersSprites, _ := utils.GetRoostersSprites("https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/sprites.json")

	roosterIndex := prototypes.FindIndex(allValidRoosters, func(rooster utils.Class) bool {
		return strings.EqualFold(rooster.Name, roosterName)
	})

	roosterImg := roostersSprites[0][roosterIndex-1]

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Galo **%s**", rooster.Name),
		Color: utils.GetRarityColor(rooster.Rarity),
		Image: &discordgo.MessageEmbedImage{
			URL: roosterImg.(string),
		},
	}

	response.Embeds = []*discordgo.MessageEmbed{embed}
	response.Components = buttons

	/*message, _ := */
	s.ChannelMessageSendComplex(m.ChannelID, response)

	// handler.CreateMessageComponentCollector(message, func(i *discordgo.Interaction) {
	// 	switch i.MessageComponentData().CustomID {
	// 	case "skins":
	// 		fmt.Println("skins", rooster)
	// 	case "skills":
	// 		fmt.Println("skills", rooster)
	// 	}
	// }, 0)
}

func ShowRoosterSkins(s *discordgo.Session, i *discordgo.Interaction, rooster *utils.Class) {

}
