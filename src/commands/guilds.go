package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/prototypes"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:      "guilds",
		Aliases:   []string{"sv", "servidores", "svs"},
		Developer: true,
		Run:       runGuilds,
	})
}

func runGuilds(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	guilds := s.State.Guilds

	prototypes.SortSlice(guilds, func(a, b *discordgo.Guild) bool {
		return a.MemberCount < b.MemberCount
	}, true)

	var guildsFormatted []string
	for _, guild := range guilds {
		guildsFormatted = append(guildsFormatted, formatGuildInfo(guild))
	}

	const pageSize = 5
	totalPages := (len(guildsFormatted) + pageSize - 1) / pageSize
	page := 0

	createDescription := func(page int) string {
		start := page * pageSize
		end := start + pageSize
		if end > len(guildsFormatted) {
			end = len(guildsFormatted)
		}

		var sb strings.Builder
		for _, str := range guildsFormatted[start:end] {
			sb.WriteString(str)
			sb.WriteString("\n")
		}
		return sb.String()
	}

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("**%d/%d**", page+1, totalPages),
		Description: createDescription(page),
		Color:       utils.ColorDefault,
	}

	message, _ := discord.NewMessage(s, m.ChannelID, m.ID).
		WithEmbed(embed).
		WithButton("Voltar", "prev_page", discordgo.PrimaryButton, nil).
		WithButton("Avançar", "next_page", discordgo.PrimaryButton, nil).
		Send()

	handler.CreateMessageComponentCollector(message, func(i *discordgo.Interaction) {
		switch i.MessageComponentData().CustomID {
		case "prev_page":
			if page > 0 {
				page--
			}
		case "next_page":
			if page+1 < totalPages {
				page++
			}
		}

		embed.Title = fmt.Sprintf("**%d/%d**", page+1, totalPages)
		embed.Description = createDescription(page)

		s.InteractionRespond(i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
				Components: []discordgo.MessageComponent{
					&discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							&discordgo.Button{
								Label:    "Voltar",
								Style:    discordgo.PrimaryButton,
								CustomID: "prev_page",
								Disabled: page == 0,
							},
							&discordgo.Button{
								Label:    "Avançar",
								Style:    discordgo.PrimaryButton,
								CustomID: "next_page",
								Disabled: page+1 >= totalPages,
							},
						},
					},
				},
			},
		})
	}, 0)
}

func formatGuildInfo(guild *discordgo.Guild) (text string) {
	text += fmt.Sprintf("### %s\n- ID: %s\n- Membros: %d\n- Owner ID: %s", guild.Name, guild.ID, guild.MemberCount, guild.OwnerID)
	return
}
