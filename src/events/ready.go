package events

import (
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var once sync.Once
var OwnersIDs = []string{"339314508621283329", "1030277251377410068"}

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	executeOnReadyOnce(s)

	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}

func executeOnReadyOnce(s *discordgo.Session) {
	once.Do(func() {
		log.Println("Initializing Interactions and Messages workers")
		InitInteractionWorkers(s)
		InitMessageWorkers(s)
	})
}
