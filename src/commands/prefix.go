package commands

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "prefix",
		Aliases:     []string{"prefixo"},
		Description: "Prefix.Help",
		Category:    "config",
		Usage:       "Prefix.Usage",
		Run:         runPrefix,
	})
}

func runPrefix(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	author := m.Author
	memberPermissions, err := discord.MemberHasPermission(s, m.Message, discordgo.PermissionManageServer)

	if err != nil {
		log.Fatalln("erro ao verificar permissões:", err)
		return
	}

	user := database.User.GetUser(ctx, m.Author)

	if !memberPermissions {
		response := services.Translate("Prefix.NoPermission", &user, author.Mention())
		discord.NewMessage(s, m.ChannelID, m.ID).
			WithContent(response).
			Send()
		return
	}

	if len(args) < 1 {
		response := services.Translate("Prefix.NoArgs", &user, author.Mention())
		discord.NewMessage(s, m.ChannelID, m.ID).
			WithContent(response).
			Send()
		return
	}

	newPrefix := args[0]

	if len(newPrefix) > 5 {
		response := services.Translate("Prefix.TooLong", &user, author.Mention())
		discord.NewMessage(s, m.ChannelID, m.ID).
			WithContent(response).
			Send()
		return
	}

	validPrefixRegex := regexp.MustCompile(`^[A-Za-z0-9~` + "`" + `!@#$%^&*()_+\-={}|:;<>,.?\/']+$`)

	if !validPrefixRegex.MatchString(newPrefix) {
		response := services.Translate("Prefix.Invalid", &user, author.Mention())
		discord.NewMessage(s, m.ChannelID, m.ID).
			WithContent(response).
			Send()
		return
	}

	newPrefix = strings.ToLower(newPrefix)

	database.Guild.UpdateGuild(ctx, m.GuildID, func(g schemas.Guild) schemas.Guild {
		g.Prefix = newPrefix
		response := services.Translate("Prefix.Success", &user, map[string]interface{}{
			"Prefix": newPrefix,
			"User":   author.Mention(),
		})
		discord.NewMessage(s, m.ChannelID, m.ID).
			WithContent(response).
			Send()
		return g
	})

}
