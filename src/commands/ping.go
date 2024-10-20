package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:        "ping",
		Aliases:     []string{"pong"},
		Cooldown:    5,
		Run:         runPing,
		Category:    "general",
		Description: "Ping.Description",
	})
}

func runPing(_ context.Context, s *discordgo.Session, m *discordgo.MessageCreate, _ []string) {
	latency := s.HeartbeatLatency().Milliseconds()
	shardId := s.ShardID + 1
	shardCount := s.ShardCount

	res := fmt.Sprintf("Pong! Shard (**%d**/**%d**)\nGateway Ping: **%dms**", shardId, shardCount, latency)

	response := &discordgo.MessageSend{
		Content: res,
		Reference: &discordgo.MessageReference{
			MessageID: m.ID,
		},
	}

	s.ChannelMessageSendComplex(m.ChannelID, response)
}
