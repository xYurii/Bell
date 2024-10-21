package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}
