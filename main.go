package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/servusdei2018/shards/v2"
	_ "github.com/xYurii/Bell/src/commands"
	_ "github.com/xYurii/Bell/src/components"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/events"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// connect to the database
	_, err := database.InitDatabase(database.GetEnvDatabaseConfig())
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	conn := fmt.Sprintf("Bot %s", os.Getenv("DISCORD_CLIENT_TOKEN"))
	s, _ := shards.New(conn)
	s.Intent = discordgo.IntentsAll

	s.AddHandler(events.OnReady)
	s.AddHandler(events.OnMessageCreate)
	s.AddHandler(events.OnInteractionCreate)
	s.AddHandler(events.OnPresenceUpdate)
	s.AddHandler(events.OnMemberJoin)

	if err := s.Start(); err != nil {
		log.Fatalf("Error starting shards: %v", err)
	}

	handler.ReadyAt = time.Now()
	// load the asura roosters effects and cosmetics:
	utils.GetCosmetics()
	utils.GetEffects()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigch

	if err := s.Shutdown(); err != nil {
		log.Printf("could not close session: %s", err)
	}
}
