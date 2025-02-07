package components

import (
	"context"
	"fmt"
	"os"

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
	err := handler.DeferInteraction(s, i.Interaction)
	if err != nil {
		fmt.Println(err)
	}

	res := fmt.Sprintf(
		"Para contabilizar seu tempo com o status ativo, siga estes passos:\n\n"+
			"1️⃣ **Defina esta mensagem no seu status personalizado:**\n➜ `%s`\n\n"+
			"2️⃣ **Mantenha seu perfil online!** O status deve estar como **Online, Ausente ou Não Perturbe**. "+
			"Se estiver **Invisível/Offline**, o tempo **não será contado**.\n\n"+
			"3️⃣ **Ao salvar o status, marque a opção \"Não limpar automaticamente\"**, "+
			"assim ele permanecerá ativo sem precisar ser redefinido manualmente.",
		events.TargetStatus,
	)

	videoFile, err := os.Open("status_tutorial.mp4")
	if err != nil {
		fmt.Println(err)
	}
	defer videoFile.Close()

	response := &discordgo.WebhookParams{
		Content: res,
		Flags:   discordgo.MessageFlagsEphemeral,
		Files: []*discordgo.File{
			{
				Name:   "status_tutorial.mp4",
				Reader: videoFile,
			},
		},
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, response)
	if err != nil {
		fmt.Println(err)
	}
}
