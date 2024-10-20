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
		Description: "Language.Help",
	})
}

func runLanguage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	user := database.User.GetUser(ctx, m.Author)
	placeholder := services.Translate("Language.Menuplaceholder", &user)

	options := generateOptions(&user)
	msg, _ := discord.NewMessage(s, m.ChannelID, m.ID).
		WithSelectMenu("language-change", placeholder, options, 1, 1, false).
		WithEmbed(generateEmbed(&user, m.Author)).
		Send()

	handler.CreateMessageComponentCollector(msg, func(i *discordgo.Interaction) {
		l := i.MessageComponentData().Values[0]
		database.User.UpdateUser(ctx, m.Author, func(u schemas.User) schemas.User {
			u.Language = l
			res := services.Translate("Language.Success", &schemas.User{Language: l}, l)
			handler.RespondInteraction(s, i, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
			return u
		})
	}, 0)
}

func generateEmbed(user *schemas.User, author *discordgo.User) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       services.Translate("Language.Title", user),
		Description: services.Translate("Language.Embeddescription", user, user.Language),
		Color:       utils.ColorDefault,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: author.AvatarURL("2048"),
		},
	}
}

func generateOptions(user *schemas.User) []discordgo.SelectMenuOption {
	return []discordgo.SelectMenuOption{
		{
			Label:   services.Translate("Language.Portuguese", user),
			Value:   "pt-BR",
			Default: user.Language == "pt-BR",
		},
		{
			Label:   services.Translate("Language.English", user),
			Value:   "en-US",
			Default: user.Language == "en-US",
		},
	}
}
