package events

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

func handleLocalComponent(s *discordgo.Session, i *discordgo.InteractionCreate, ctx context.Context) {
	if collector, exists := handler.GetMessageComponentCollector(i.Message); exists {
		go func() {
			mutex.Lock()
			defer mutex.Unlock()
			collector.Callback(i.Interaction)
		}()
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
		handleLocalComponent(s, i, ctx)
	}
}
