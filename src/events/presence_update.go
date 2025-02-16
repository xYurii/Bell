package events

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
)

const TargetStatus = "asura bot (rinha de galo) https://asurabot.com.br/"

var TargetsStatus = []string{"asura bot (rinha de galo) https://acnologla.github.io/asura-site/", TargetStatus}

func OnPresenceUpdate(s *discordgo.Session, evt *discordgo.PresenceUpdate) {
	if evt.User.Bot {
		return
	}

	hasCustomStatus := false
	for _, activity := range evt.Activities {
		if activity.Type == discordgo.ActivityTypeCustom && activity.State != "" {
			for _, target := range TargetsStatus {
				if strings.Contains(activity.State, target) {
					hasCustomStatus = true
					break
				}
			}
		}
		if hasCustomStatus {
			break
		}
	}

	if hasCustomStatus {
		startTracking(evt.User)
	} else {
		stopTracking(evt.User)
	}
}

func stopTracking(user *discordgo.User) {
	if startTime, exists := handler.UserStatusTracking[user.ID]; exists {
		duration := time.Now().Unix() - startTime
		fmt.Printf("Usuário %s usou o status '%s' por %d segundos.\n", user.ID, TargetStatus, duration)
		delete(handler.UserStatusTracking, user.ID)
		database.User.UpdateUser(context.Background(), user, func(u schemas.User) schemas.User {
			u.StatusTime += duration
			return u
		})
	}
}

func startTracking(user *discordgo.User) {
	if _, exists := handler.UserStatusTracking[user.ID]; !exists {
		handler.UserStatusTracking[user.ID] = time.Now().Unix()
		fmt.Printf("Usuário %s começou a usar o status: %s\n", user.ID, TargetStatus)
	}
}
