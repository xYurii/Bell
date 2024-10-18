package events

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var mutex sync.RWMutex

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}
