package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:     "xp",
		Aliases:  []string{"xp-hierarchy", "xphierarchy", "level", "lvl"},
		Cooldown: 5,
		Run:      runXp,
	})
}

func runXp(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	var res string
	var totalNeededXp int

	for i := 1; i < 50; i++ {
		currXp := utils.CalcXP(i)
		nextXp := utils.CalcXP(i + 1)
		totalNeededXp += nextXp - currXp

		res += fmt.Sprintf("Nível %d -> %d: %d XP\n", i, i+1, nextXp-currXp)
	}

	res += fmt.Sprintf("Total de XP até o nível 50: %d XP\n", totalNeededXp)

	response := &discordgo.MessageSend{
		Content: res,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, response)
}
