package events

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/utils/tasks"
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

		// start raffles tasks
		rm := tasks.NewRaffleManager(database.Database, s)
		err := rm.EnsureRafflesExist(context.Background())
		if err != nil {
			log.Fatalf("Error ensuring raffles exist: %s", err)
		}

		go startRaffleExpirationTicker(rm)
	})
}

func startRaffleExpirationTicker(rm *tasks.RaffleManager) {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		rm.EndExpiredRaffles()
	}
}
