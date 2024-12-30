package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:      "set_time",
		Aliases:   []string{"setartempo", "settime", "addtime", "addtempo"},
		Cooldown:  5,
		Run:       runSetTime,
		Developer: true,
	})
}

func runSetTime(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	reply := discord.NewMessage(s, m.ChannelID, m.ID)
	if len(args) < 2 {
		reply.WithContent("Você precisa informar o ID do usuário e o tempo a ser adicionado.\nExemplo: j!settime <@user> <tempo> [Tempo: 1d, 1h...]").Send()
		return
	}

	user := discord.GetUser(s, m, args)
	if user.Bot {
		reply.WithContent("Você não pode adicionar tempo para um bot.").Send()
	}
	seconds := utils.ConvertStringToSeconds(strings.Join(args[1:], " "))

	database.User.UpdateUser(ctx, user, func(u schemas.User) schemas.User {
		u.StatusTime = int64(seconds)
		return u
	})
	newTime := utils.FormatDuration(int64(seconds))
	reply.WithContent(fmt.Sprintf("**%s** agora possui **%s**", user.Username, newTime)).Send()
}
