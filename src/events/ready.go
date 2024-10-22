package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var OwnersIDs = []string{"339314508621283329", "1030277251377410068"}

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}
