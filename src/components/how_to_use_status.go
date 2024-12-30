package components

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/events"
	"github.com/xYurii/Bell/src/handler"
)

func init() {
	handler.RegisterComponent(handler.Component{
		Name: "how-to-start-status-time",
		Type: discordgo.ButtonComponent,
		Run:  howToUseStatus,
	})
}

func howToUseStatus(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	res := fmt.Sprintf("Para iniciar a contagem do tempo com o status ativo, basta usar a mensagem `%s` no seu status customizado e manter o perfil em qualquer modo exceto o invisível/offline! Lembre-se de marcar a opção \"não limpar\" na hora de salvar o status.", events.TargetStatus)
	handler.RespondInteraction(s, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
}
