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
	if len(args) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Você não informou o nome do galo que deseja visualizar!")
		return
	}

	roosterName := args[0]

	if roosterName == "" {
		s.ChannelMessageSend(m.ChannelID, "Você não informou o nome do galo que deseja visualizar!")
		return
	}

	roosters, err := utils.GetJSON("https://raw.githubusercontent.com/Acnologla/asura/master/resources/galo/class.json")

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Erro ao buscar a lista de galo!")
		return
	}

	allRoostersValidNames := utils.GetRoostersNames(roosters)
	roosterExists := prototypes.Includes(allRoostersValidNames, func(name string) bool {
		return strings.EqualFold(name, roosterName)
	})

	if !roosterExists {
		s.ChannelMessageSend(m.ChannelID, "O galo informado não existe!")
		return
	}

	response := &discordgo.MessageSend{
		Content: fmt.Sprintf("Galo **%s** existe!", roosterName),
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, response)
}
