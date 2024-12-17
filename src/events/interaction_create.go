package events

import (
	"context"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/xYurii/Bell/src/handler"
)

const Workers = 128

var WorkersArray = make([]bool, Workers)
var InteractionCreateChannel = make(chan *discordgo.InteractionCreate, Workers)

var ComponentLock = sync.RWMutex{}

func HandleComponent(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate) {
	ComponentLock.RLock()
	defer ComponentLock.RUnlock()

	collector, existsCollector := handler.GetMessageComponentCollector(i.Message)
	if existsCollector {
		collector.Lock()
		defer collector.Unlock()
		collector.Callback(i.Interaction)
		return
	}

	globalComponent, existsGlobal := handler.GetComponent(i.MessageComponentData().CustomID)
	if existsGlobal {
		globalComponent.Run(ctx, s, i)
		return
	}

	res := "O cache desta interação expirou! Use o respectivo comando novamente."
	handler.RespondInteraction(s, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, res, discordgo.MessageFlagsEphemeral)
}

func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		InteractionCreateChannel <- i
	}
}

func Worker(id int, s *discordgo.Session) {
	for itc := range InteractionCreateChannel {
		WorkersArray[id] = true
		ctx := context.Background()
		HandleComponent(ctx, s, itc)
		WorkersArray[id] = false
	}
}

func InitInteractionWorkers(s *discordgo.Session) {
	for i := 0; i < Workers; i++ {
		go Worker(i, s)
	}
}

func GetFreeWorkers() int {
	var freeWorkers int
	for _, worker := range WorkersArray {
		if !worker {
			freeWorkers++
		}
	}
	return freeWorkers
}
