package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "raffle",
		Aliases:     []string{"rifa", "rf"},
		Description: "Raffle.Help",
		Cooldown:    15,
		Run:         runRaffle,
	})
}

func runRaffle(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	user := database.User.GetUser(ctx, m.Author)
	if len(args) < 1 {
		res := services.Translate("Raffle.NoArgs", &user, m.Author.Mention())
		discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
	}
}
