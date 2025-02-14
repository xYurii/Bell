package events

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func OnMemberJoin(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	guildId := os.Getenv("SUPPORT_GUILD")
	if m.GuildID != guildId {
		return
	}

	channel, err := s.UserChannelCreate(m.User.ID)
	if err != nil {
		return
	}

	guildData := database.Guild.GetGuild(context.Background(), m.GuildID)

	description := fmt.Sprintf(
		"Sabia que você pode ganhar recompensas ao divulgar o bot no seu status?"+
			" Para saber como, use o comando `%sstatus` no servidor do Asura!",
		guildData.Prefix,
	)

	embed := discord.NewEmbedBuilder().
		WithDescription(description).
		WithField("**Recompensas**", "A cada **10 dias** com o status ativo, você recebe **10.000** moedas;\nA cada **40 dias** com o status ativo, você ganha uma lootbox lendária!", false).
		WithColor(utils.ColorDefault)

	embed2 := discord.NewEmbedBuilder().
		WithDescription("Gostaria de nos informar por onde você conheceu o bot? Para isso, basta escolher a opção referente nos botões abaixo.").
		WithColor(utils.ColorDefault)

	reply := discord.NewMessage(s, channel.ID, "").
		WithEmbed(embed.Build()).
		WithEmbed(embed2.Build())

	buttonsIds := []string{"Amigo", "Através do status", "Vi no site", "Outros", "Algum servidor"}
	for _, originalId := range buttonsIds {
		id := "guild_" + strings.ReplaceAll(originalId, " ", "-")
		reply.WithButton(originalId, id, discordgo.PrimaryButton, nil)
	}

	reply.Send()

	// _, err = s.ChannelMessageSend(channel.ID, "oi")
	// if err != nil {
	// 	fmt.Println("não deu p enviar dm.", err.Error())
	// }
}
