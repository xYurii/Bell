package handler

import "github.com/bwmarrin/discordgo"

func RespondInteraction(s *discordgo.Session, i *discordgo.Interaction, responseType discordgo.InteractionResponseType, content string, flags ...discordgo.MessageFlags) error {
	response := &discordgo.InteractionResponse{
		Type: responseType,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	}

	if len(flags) > 0 {
		response.Data.Flags = flags[0]
	}

	return s.InteractionRespond(i, response)
}

func DeferInteraction(s *discordgo.Session, i *discordgo.Interaction) error {
	return s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}
