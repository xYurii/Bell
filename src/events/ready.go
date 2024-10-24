package events

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var once sync.Once
var OwnersIDs = []string{"339314508621283329", "1030277251377410068"}

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	once.Do(func() {
		log.Println("Initializing workers")
		InitInteractionWorkers(s)
		InitMessageWorkers(s)
	})

	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}
