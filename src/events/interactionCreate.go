package events

import (
	"context"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

var mutex sync.RWMutex

func handleLocalComponent(s *discordgo.Session, i *discordgo.InteractionCreate, ctx context.Context) {
	if collector, exists := handler.GetMessageComponentCollector(i.Message); exists {
		collector.Callback(i.Interaction)
	} else {
		handleGlobalComponent(s, i, ctx)
	}
}

func handleGlobalComponent(s *discordgo.Session, i *discordgo.InteractionCreate, ctx context.Context) {
	globalComponent, existsGlobal := handler.GetComponent(i.MessageComponentData().CustomID)
	if existsGlobal {
		globalComponent.Run(ctx, s, i)
	} else {
		res := "O cache desta interação expirou! Use o respectivo comando novamente."
		handler.RespondInteraction(s, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
	}
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := context.Background()
	if i.Type == discordgo.InteractionMessageComponent {
		go handleLocalComponent(s, i, ctx)
	}
}
