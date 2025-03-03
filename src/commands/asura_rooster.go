package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/prototypes"
	"github.com/xYurii/Bell/src/utils"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "galo",
		Aliases:     []string{"rooster"},
		Cooldown:    5,
		Run:         runAsuraRooster,
		Category:    "asura",
		Description: "Asura.Rooster",
	})
}

func runAsuraRooster(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	response := &discordgo.MessageSend{
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	}

	buttons := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Habilidades",
					Style:    discordgo.SecondaryButton,
					CustomID: "skills",
				},
				discordgo.Button{
					Label:    "Skins",
					Style:    discordgo.SuccessButton,
					CustomID: "skins",
				},
			},
		},
	}

	buttonsWithValidRoostersNames := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Nome de todos os Galos",
					Style:    discordgo.PrimaryButton,
					CustomID: "allRoostersNames",
				},
			},
		},
	}

	if len(args) == 0 {
		response.Content = "Você não informou o nome do galo que deseja visualizar!"
		response.Components = buttonsWithValidRoostersNames
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	roosterName := args[0]
	resets := 0

	if len(args) > 1 {
		var err error
		resets, err = strconv.Atoi(args[1])
		if err != nil {
			resets = 0
		}
		if resets < 0 {
			resets = 0
		}
	}

	if roosterName == "" {
		response.Content = "Você deve informar o nome do galo que deseja visualizar!"
		response.Components = buttonsWithValidRoostersNames
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	roosters, err := utils.GetRoostersClasses("https://raw.githubusercontent.com/Acnologla/asura-site/main/public/resources/class.json")

	if err != nil {
		response.Content = "Eu não consegui buscar a lista dos nomes dos galos... Tente novamente."

		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	allValidRoosters := utils.GetRoostersNames(roosters)
	roosterExists := prototypes.Includes(allValidRoosters, func(roosterData utils.Class) bool {
		return strings.EqualFold(roosterData.Name, roosterName)
	})

	if !roosterExists {
		response.Content = fmt.Sprintf("O galo **%s** não existe!", roosterName)
		response.Components = buttonsWithValidRoostersNames
		s.ChannelMessageSendComplex(m.ChannelID, response)
		return
	}

	rooster := prototypes.Find(allValidRoosters, func(rooster utils.Class) bool {
		return strings.EqualFold(rooster.Name, roosterName)
	})

	roostersSprites, _ := utils.GetRoostersSprites("https://raw.githubusercontent.com/Acnologla/asura-site/main/public/resources/sprites.json")

	roosterIndex := prototypes.FindIndex(allValidRoosters, func(rooster utils.Class) bool {
		return strings.EqualFold(rooster.Name, roosterName)
	})

	roosterImg := roostersSprites[0][roosterIndex-1]

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Galo **%s** - **%s**", rooster.Name, rooster.Rarity.String()),
		Color: rooster.Rarity.Color(),
		Image: &discordgo.MessageEmbedImage{
			URL: roosterImg,
		},
	}

	response.Embeds = []*discordgo.MessageEmbed{embed}
	response.Components = buttons

	message, _ := s.ChannelMessageSendComplex(m.ChannelID, response)

	handler.CreateMessageComponentCollector(message, func(i *discordgo.Interaction) {
		switch i.MessageComponentData().CustomID {
		case "skins":
			showRoosterSkins(s, i, roosterIndex)
		case "skills":
			skills := showRoosterSkills(&rooster, float64(resets))

			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("%s - Resets: %d", rooster.Name, resets),
				Color:       rooster.Rarity.Color(),
				Description: skills,
			}

			response := &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			}
			s.InteractionRespond(i, response)
		}
	}, 1*time.Minute)
}

func showRoosterSkills(rooster *utils.Class, resets float64) string {
	skills, _ := utils.GetRoosterSkills(rooster)
	skillsMap := prototypes.Map(skills, func(skill *utils.Skill) string {
		min, max := utils.CalcDamage(skill.Damage[0], skill.Damage[1], resets)
		level := strconv.Itoa(skill.Level)
		if skill.Evolved {
			level = fmt.Sprintf("⭐ %d", skill.Level)
		}
		// text := fmt.Sprintf("%s [**%s**]: **%d** - **%d**", skill.Name, level, min, max)
		text := fmt.Sprintf("[**%s**] **%s** ➜ **%d** - **%d**", level, skill.Name, min, max)

		if skill.Effect[0] != 0 || skill.Effect[1] != 0 {
			effect := utils.Effects[int(skill.Effect[1])]
			minEffect, maxEffect := utils.CalcDamage(effect.Range[0], effect.Range[1], resets)
			turns := effect.Turns
			turnsText := ""
			self := ""
			if turns > 0 {
				turnsText = fmt.Sprintf(" (%d turnos)", turns)
			}
			if effect.Self {
				self = "em si mesmo"
			}
			text += fmt.Sprintf("\nTem %d%% de Chance de causar **%s** %s [**%d** - **%d**]%s", int(skill.Effect[0]*100), effect.Name, self, minEffect, maxEffect, turnsText)
		}

		return text
	})

	var text string
	for _, skill := range skillsMap {
		text += skill + "\n\n"
	}

	return text
}

func showRoosterSkins(s *discordgo.Session, i *discordgo.Interaction, roosterIndex int) {
	var embeds []*discordgo.MessageEmbed

	for _, c := range utils.Cosmetics {
		if c.Extra == roosterIndex {
			embed := &discordgo.MessageEmbed{
				Title: fmt.Sprintf("**%s** - **%s**", c.Name, c.Rarity.String()),
				Color: c.Rarity.Color(),
				Image: &discordgo.MessageEmbedImage{
					URL: c.Value,
				},
			}
			embeds = append(embeds, embed)
		}
	}

	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Nenhuma skin encontrada.",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}

	if len(embeds) > 0 {
		response.Data = &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		}
	}

	s.InteractionRespond(i, response)
}
