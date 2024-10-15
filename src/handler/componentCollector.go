package handler

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Collector struct {
	Callback func(*discordgo.Interaction)
	Timeout  *time.Time
}

var Collectors = map[string]Collector{}
var mu sync.Mutex

func CreateMessageComponentCollector(msg *discordgo.Message, callback func(*discordgo.Interaction), timeout time.Duration, onExpire ...func()) {
	mu.Lock()
	defer mu.Unlock()

	var expires *time.Time
	if timeout > 0 {
		expTime := time.Now().Add(timeout)
		expires = &expTime
	}

	Collectors[msg.ID] = Collector{
		Callback: callback,
		Timeout:  expires,
	}

	if expires != nil {
		go func() {
			time.Sleep(timeout)
			mu.Lock()
			defer mu.Unlock()
			onExpire[0]()
			delete(Collectors, msg.ID)
		}()
	}
}

func DeleteComponentCollector(msg *discordgo.Message) {
	mu.Lock()
	defer mu.Unlock()

	delete(Collectors, msg.ID)
}

func GetMessageComponentCollector(msg *discordgo.Message) (*Collector, bool) {
	mu.Lock()
	defer mu.Unlock()

	collector, exists := Collectors[msg.ID]
	return &collector, exists
}

func CleanComponentCollectors() {
	mu.Lock()
	defer mu.Unlock()

	for k, v := range Collectors {
		if v.Timeout != nil && time.Now().After(*v.Timeout) {
			delete(Collectors, k)
		}
	}
}
