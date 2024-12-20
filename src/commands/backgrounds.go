package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "backgrounds",
		Aliases:     []string{"bgs", "bg"},
		Cooldown:    5,
		Run:         runBackgrounds,
		Category:    "asura",
		Description: "Asura.Backgrounds",
	})
}

func runBackgrounds(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	page := 0
	backgrounds := utils.GetBackgrounds()
	background := backgrounds[page]

	buttons := CreateBackgroundsButtons(page, m.Author.ID)

	reply := discord.NewMessage(s, m.ChannelID, m.ID).
		WithEmbed(BuildBackgroundsEmbed(background, page, len(backgrounds))).
		WithButtons(buttons)

		// response := &discordgo.MessageSend{
		// 	Embed:      BuildBackgroundsEmbed(background, page, len(backgrounds)),
		// 	Components: CreateBackgroundsButtons(page, m.Author.ID),
		// 	Reference: &discordgo.MessageReference{
		// 		MessageID: m.ID,
		// 	},
		// }

		/* _, err := */
	reply.Send() // s.ChannelMessageSendComplex(m.ChannelID, response)
	// if err != nil {
	// 	s.ChannelMessageSend(m.ChannelID, "Erro ao enviar o embed.")
	// 	return
	// }
}
func BuildBackgroundsEmbed(background *utils.Cosmetic, page, backgroundsLen int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("**%s** - **%s**", background.Name, background.Rarity.String()),
		Image: &discordgo.MessageEmbedImage{
			URL: background.Value,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Página %d/%d", page+1, backgroundsLen),
		},
		Color: background.Rarity.Color(),
	}
}

func CreateBackgroundsButtons(page int, userID string) []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					CustomID: "backgrounds_first_" + strconv.Itoa(page) + "_" + userID,
					Label:    "Home",
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					CustomID: "backgrounds_previous_" + strconv.Itoa(page) + "_" + userID,
					Emoji:    &discordgo.ComponentEmoji{Name: "DoubleLeftArrow", ID: "1272032089507762242"},
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					CustomID: "backgrounds_next_" + strconv.Itoa(page) + "_" + userID,
					Emoji:    &discordgo.ComponentEmoji{Name: "DoubleRightArrow", ID: "1272031913888059468"},
				},
				discordgo.Button{
					Style:    discordgo.PrimaryButton,
					CustomID: "backgrounds_last_" + strconv.Itoa(page) + "_" + userID,
					Label:    "End",
				},
			},
		},
	}
}
