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

	// fmt.Println(services.Translate("Help.Title", &schemas.User{Language: "pt-BR"}))

	// load the asura roosters effects and cosmetics:
	utils.GetCosmetics()
	utils.GetEffects()

	// connect to the database
	database.InitDatabase(database.GetEnvDatabaseConfig())

	conn := fmt.Sprintf("Bot %s", os.Getenv("DISCORD_CLIENT_TOKEN"))
	s, _ := shards.New(conn)
	s.Intent = discordgo.IntentsAll

	s.AddHandler(events.OnReady)
	s.AddHandler(events.OnMessageCreate)
	s.AddHandler(events.OnInteractionCreate)

	if err := s.Start(); err != nil {
		log.Fatalf("Error starting shards: %v", err)
	}

	events.InitInteractionWorkers(s.Gateway)
	handler.ReadyAt = time.Now()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigch

	if err := s.Shutdown(); err != nil {
		log.Printf("could not close session: %s", err)
	}
}
