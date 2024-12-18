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
	if len(args) > 0 {
		if args[0] == "top" || args[0] == "rank" {
			showTopUserStatusTime(ctx, s, m)
			return
		}
	}

	user := discord.GetUser(s, m, args)
	data := database.User.GetUser(ctx, user)
	formattedTime := utils.FormatDuration(data.StatusTime)

	description := fmt.Sprintf("Tempo total com o status ativo: **%s**", formattedTime)

	if duration, exists := handler.UserStatusTracking[user.ID]; exists {
		cacheDuration := utils.FormatDuration(time.Now().Unix() - duration)
		description += fmt.Sprintf("\nTempo em cache com o status ativo: **%s**", cacheDuration)
	}

	embed := discord.NewEmbedBuilder().
		WithColor(utils.ColorDefault).
		WithDescription(description).
		WithFooter("O tempo em cache é somado ao tempo total a cada ~30s", user.AvatarURL("256")).
		WithTimestamp(time.Now().Format(time.RFC3339))

	discord.NewMessage(s, m.ChannelID, m.ID).WithEmbed(embed.Build()).Send()
}

func showTopUserStatusTime(ctx context.Context, s *discordgo.Session, evt *discordgo.MessageCreate) {
	users := database.User.SortUsers(ctx, 10, "status_time")
	var text string
	for i, u := range users {
		member, _ := s.State.Member(evt.GuildID, u.ID)
		username := u.ID
		if member != nil {
			username = member.User.Username
		} else {
			user, _ := s.User(u.ID)
			if user != nil {
				username = user.Username
			}
		}
		formattedTime := utils.FormatDuration(u.StatusTime)
		text += fmt.Sprintf("**[%d].** `%s`: **%s**\n", i+1, username, formattedTime)
	}

	guild, _ := s.State.Guild(evt.GuildID)

	embed := discord.NewEmbedBuilder().
		WithColor(utils.ColorDefault).
		WithDescription(text).
		WithFooter(guild.Name, guild.IconURL("256"))

	discord.NewMessage(s, evt.ChannelID, evt.ID).WithEmbed(embed.Build()).Send()
}
