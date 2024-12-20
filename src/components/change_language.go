package components

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
)

func init() {
	handler.RegisterComponent(handler.Component{
		Name: "language-change",
		Type: discordgo.ButtonComponent,
		Run:  changeLanguage,
	})
}

func changeLanguage(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	lang := i.MessageComponentData().Values[0]

	database.User.UpdateUser(ctx, i.Member.User, func(u schemas.User) schemas.User {
		u.Language = lang
		return u
	})

	res := services.Translate("Language.Success", &schemas.User{Language: lang}, lang)
	handler.RespondInteraction(s, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
}
