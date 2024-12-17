package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:     "status",
		Cooldown: 5,
		Run:      runStatus,
		Category: "general",
	})
}

func runStatus(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	data := database.User.GetUser(ctx, m.Author)
	formattedTime := formatDuration(data.StatusTime)

	embed := &discordgo.MessageEmbed{
		Title:       "Tempo com o status",
		Description: fmt.Sprintf("Você esteve com o status durate: **%s**.", formattedTime),
		Color:       utils.ColorDefault,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL: m.Author.AvatarURL("256"),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if duration, exists := handler.UserStatusTracking[m.Author.ID]; exists {
		cacheDuration := formatDuration(time.Now().Unix() - duration)
		embed.Description += fmt.Sprintf("\nVocê possui **%s** de tempo acumulado por estar com o status ativo!", cacheDuration)
	}

	discord.NewMessage(s, m.ChannelID, m.ID).WithEmbed(embed).Send()
}

func formatDuration(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	var result string
	if days > 0 {
		result += fmt.Sprintf("%dd, ", days)
	}
	if hours > 0 {
		result += fmt.Sprintf("%dh, ", hours)
	}
	if minutes > 0 {
		result += fmt.Sprintf("%dm, ", minutes)
	}
	result += fmt.Sprintf("%ds", remainingSeconds)

	return result
}
