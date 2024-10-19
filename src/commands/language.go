package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "language",
		Aliases:     []string{"lang", "idioma", "linguagem"},
		Cooldown:    5,
		Run:         runLanguage,
		Category:    "general",
		Description: "language.help",
	})
}

func runLanguage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	user := database.User.GetUser(ctx, m.Author.ID)

	options := generateOptions(&user)
	discord.NewMessage(s, m.ChannelID, m.ID).
		WithSelectMenu("language-change", "Escolha o idioma que você quer que eu fale com você!", options, 1, 1, false).
		WithEmbed(generateEmbed(&user, m.Author)).
		Send()
}

func generateEmbed(user *schemas.User, author *discordgo.User) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       services.Translate("language.title", user),
		Description: services.Translate("language.embeddescription", user, user.Language),
		Color:       utils.ColorDefault,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: author.AvatarURL("2048"),
		},
	}
}

func generateOptions(user *schemas.User) []discordgo.SelectMenuOption {
	return []discordgo.SelectMenuOption{
		{
			Label:   services.Translate("language.portuguese", user),
			Value:   "pt",
			Default: user.Language == "pt-BR",
		},
		{
			Label:   services.Translate("language.english", user),
			Value:   "en",
			Default: user.Language == "en-US",
		},
	}
}
