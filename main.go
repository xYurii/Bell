package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/servusdei2018/shards/v2"
	_ "github.com/xYurii/Bell/src/commands"
	_ "github.com/xYurii/Bell/src/components"
	"github.com/xYurii/Bell/src/events"
	"github.com/xYurii/Bell/src/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// load the asura roosters effects and cosmetics:
	utils.GetCosmetics()
	utils.GetEffects()

	conn := fmt.Sprintf("Bot %s", os.Getenv("DISCORD_CLIENT_TOKEN"))
	s, _ := shards.New(conn)
	s.Intent = discordgo.IntentsGuildMessages

	s.AddHandler(events.OnReady)
	s.AddHandler(events.OnMessageCreate)
	s.AddHandler(events.OnInteractionCreate)

	if err := s.Start(); err != nil {
		log.Fatalf("Error starting shards: %v", err)
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigch

	if err := s.Shutdown(); err != nil {
		log.Printf("could not close session: %s", err)
	}
}
