package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/prototypes"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:      "guilds",
		Aliases:   []string{"sv", "servidores", "svs"},
		Developer: true,
		Run:       runGuilds,
	})
}

func runGuilds(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	guilds := s.State.Guilds
	prototypes.SortSlice(guilds, func(a, b *discordgo.Guild) bool {
		return a.MemberCount < b.MemberCount
	}, true)

	res := ""
	for _, guild := range guilds {
		res += fmt.Sprintf("Guild Name: %s\nGuild ID: \"%s\"\nMembers: %d\nOwner ID: \"%s\"\n\n", guild.Name, guild.ID, guild.MemberCount, guild.OwnerID)
	}
	res = "```yaml\n" + res + "```"
	discord.NewMessage(s, m.ChannelID, m.ID).WithContent(res).Send()
}
