package events

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
)

var (
	once      sync.Once
	OwnersIDs = []string{"339314508621283329", "1030277251377410068"}
	mu        sync.Mutex
)

func OnReady(s *discordgo.Session, r *discordgo.Ready) {
	executeOnReadyOnce(s)

	log.Printf("[SHARD: %v] logged in as %s", s.ShardID, r.User.String())
}

func executeOnReadyOnce(s *discordgo.Session) {
	once.Do(func() {
		startTrackingActiveUsers(s)
		log.Println("Initializing Interactions and Messages workers")
		InitInteractionWorkers(s)
		InitMessageWorkers(s)

		go saveUserStatus()
	})
}

func saveUserStatus() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mu.Lock()
		for id, startTime := range handler.UserStatusTracking {
			duration := time.Now().Unix() - startTime
			handler.UserStatusTracking[id] = time.Now().Unix()

			err := database.User.UpdateUser(context.Background(), &discordgo.User{ID: id}, func(u schemas.User) schemas.User {
				u.StatusTime += duration
				return u
			})
			if err != nil {
				log.Printf("Erro ao atualizar usuário %s: %v", id, err)
			}
		}
		mu.Unlock()
	}
}

func startTrackingActiveUsers(s *discordgo.Session) {
	guildId := os.Getenv("SUPPORT_GUILD")
	members, err := fetchAllGuildMembers(s, guildId)
	if err != nil {
		log.Printf("Erro ao buscar membros do servidor: %v", err)
		return
	}

	presences, err := fetchGuildPresences(s, guildId)
	if err != nil {
		log.Printf("Erro ao buscar presenças do servidor: %v", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for _, member := range members {
		if member.User.Bot {
			continue
		}

		presence, exists := presences[member.User.ID]
		if !exists {
			continue
		}

		for _, activity := range presence.Activities {
			if activity.Type == discordgo.ActivityTypeCustom && strings.Contains(activity.State, targetStatus) {
				if _, exists := handler.UserStatusTracking[member.User.ID]; !exists {
					handler.UserStatusTracking[member.User.ID] = time.Now().Unix()
				}
				break
			}
		}
	}
	log.Printf("Status tracking initialized with %d active users.", len(handler.UserStatusTracking))

}

func fetchGuildPresences(s *discordgo.Session, guildID string) (map[string]*discordgo.Presence, error) {
	presences := make(map[string]*discordgo.Presence)
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	for _, presence := range guild.Presences {
		presences[presence.User.ID] = presence
	}
	return presences, nil
}

func fetchAllGuildMembers(s *discordgo.Session, guildID string) ([]*discordgo.Member, error) {
	var allMembers []*discordgo.Member
	after := ""
	for {
		members, err := s.GuildMembers(guildID, after, 1000)
		if err != nil {
			return nil, err
		}

		allMembers = append(allMembers, members...)

		if len(members) < 1000 {
			break
		}

		after = members[len(members)-1].User.ID
	}
	return allMembers, nil
}
