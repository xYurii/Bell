package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/services"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "help",
		Aliases:     []string{"h", "ajuda"},
		Cooldown:    5,
		Run:         runHelp,
		Category:    "general",
		Description: "help.description",
	})
}

func runHelp(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	commands := map[string]map[string]bool{}

	user := database.User.GetUser(ctx, m.Author)
	guild := database.Guild.GetGuild(ctx, m.GuildID)

	if len(args) > 0 {
		cmd, exists := handler.GetCommand(args[0])
		if exists {
			showAboutCommand(s, m, cmd, &user)
			return
		}
	}

	for _, command := range handler.Commands {
		if _, ok := commands[command.Category]; !ok {
			commands[command.Category] = map[string]bool{}
		}

		commands[command.Category][command.Name] = true
	}

	embed := &discordgo.MessageEmbed{
		Title: services.Translate("Help.Title", &user),
		Color: utils.ColorDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    services.Translate("Help.Footer", &user, guild.Prefix),
			IconURL: m.Author.AvatarURL("2048"),
		},
	}

	for category, cmdsMap := range commands {
		var commandsNames []string

		for cmdName := range cmdsMap {
			commandsNames = append(commandsNames, fmt.Sprintf("`%s`", cmdName))
		}

		commandsList := strings.Join(commandsNames, ", ")
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   cases.Title(language.Portuguese).String(category),
			Value:  commandsList,
			Inline: false,
		})
	}

	discord.NewMessage(s, m.ChannelID, m.ID).
		WithEmbed(embed).
		Send()
}

func showAboutCommand(s *discordgo.Session, m *discordgo.MessageCreate, command handler.Command, user *schemas.User) {
	embed := &discordgo.MessageEmbed{
		Title: services.Translate("Aboutcommand.Title", user, command.Name),
		Color: utils.ColorDefault,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    s.State.User.Username,
			IconURL: s.State.User.AvatarURL("2048"),
		},
	}

	if command.Category != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   services.Translate("Help.Categoryfield", user),
			Value:  cases.Title(language.Portuguese).String(command.Category),
			Inline: false,
		})
	}

	if len(command.Aliases) > 0 {
		aliases := strings.Join(command.Aliases, ", ")
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   services.Translate("Help.Aliasesfield", user),
			Value:  aliases,
			Inline: false,
		})
	}

	if command.Description != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   services.Translate("Help.Descriptionfield", user),
			Value:  services.Translate(command.Description, user),
			Inline: false,
		})
	}
	if command.Usage != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   services.Translate("Help.Usagefield", user),
			Value:  services.Translate(command.Usage, user),
			Inline: false,
		})
	}

	if command.Cooldown > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   services.Translate("Help.Cooldownfield", user),
			Value:  services.Translate("Help.Cooldownvalue", user, command.Cooldown),
			Inline: false,
		})
	}

	discord.NewMessage(s, m.ChannelID, m.ID).
		WithEmbed(embed).
		Send()
}
