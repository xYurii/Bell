package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/adapter"
	"github.com/xYurii/Bell/src/database/schemas"
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
	guild := database.Guild.GetGuild(ctx, m.GuildID)
	user := database.User.GetUser(ctx, m.Author)

	if len(args) < 1 {
		res := services.Translate("Raffle.NoArgs", &user, map[string]interface{}{
			"User":   m.Author.Mention(),
			"Prefix": guild.Prefix,
		})
		discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
	}

	switch strings.ToLower(args[0]) {
	case "status":
		if len(args) < 2 {
			res := services.Translate("Raffle.NoArgs", &user, map[string]interface{}{
				"User":   m.Author.Mention(),
				"Prefix": guild.Prefix,
			})
			discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
			return
		}

		raffleType := getRaffleType(args[1])
		status, err := database.Raffle.GetRaffleStatus(ctx, raffleType)
		if err != nil {
			discord.NewMessage(s, m.ChannelID, m.ID).WithContent(err.Error()).Send()
			return
		}

		discord.NewMessage(s, m.ChannelID, m.ID).WithContent(formatRaffleTextStatus(s, &status)).Send()
	}

}

func formatRaffleTextStatus(s *discordgo.Session, status *adapter.ActualRaffleStatus) string {
	endsAt := status.Raffle.EndsAt.Unix()
	endsAtDate := fmt.Sprintf("<t:%d:F> (<t:%d:R>)", endsAt, endsAt)
	winnerName := "Unknown"
	winnerID := "Unknown"
	if status.PreviousRaffle.WinnerTicket != nil {
		winnerID = status.PreviousRaffle.WinnerTicket.UserID
		winnerUser, err := s.User(status.PreviousRaffle.WinnerTicket.UserID)
		if err == nil {
			winnerName = winnerUser.Username
		}
	}

	return fmt.Sprintf("Raffle Status:\n\nType: %s\nStarted at: <t:%d:F> (<t:%d:R>)\nEnds at: %s\nWinner: %s (%s)\nActual Price: %d\nUsers: %d\nTickets: %d\n", status.Raffle.RaffleType.String(), status.Raffle.StartedAt.Unix(), status.Raffle.StartedAt.Unix(), endsAtDate, winnerName, winnerID, status.ActualReward, status.UsersCount, status.TicketsCount)
}

func getRaffleType(arg string) schemas.RaffleType {
	var raffleType schemas.RaffleType
	switch strings.ToLower(arg) {
	case "daily":
		raffleType = schemas.DailyRaffle
	case "lightning":
		raffleType = schemas.Lightning
	default:
		raffleType = schemas.NormalRaffle

	}
	return raffleType
}
