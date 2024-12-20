package handler

import (
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Collector struct {
	Callback   func(*discordgo.Interaction)
	Timeout    *time.Time
	deleteChan chan bool
	sync.Mutex
}

var Collectors = map[string]*Collector{}
var collectorsLock = sync.RWMutex{}

func CreateMessageComponentCollector(msg *discordgo.Message, callback func(*discordgo.Interaction), timeout time.Duration, onExpire ...func()) {
	collector := &Collector{
		Callback:   callback,
		deleteChan: make(chan bool),
	}

	if timeout > 0 {
		expTime := time.Now().Add(timeout)
		collector.Timeout = &expTime
	}

	collectorsLock.Lock()
	Collectors[msg.ID] = collector
	collectorsLock.Unlock()

	go func() {
		if timeout > 0 {
			timeChannel := time.After(timeout)
			select {
			case <-collector.deleteChan:
			case <-timeChannel:
				collectorsLock.Lock()
				delete(Collectors, msg.ID)
				collectorsLock.Unlock()
				if len(onExpire) > 0 {
					onExpire[0]()
				}
			}
		} else {
			<-collector.deleteChan
		}
	}()
}

func DeleteComponentCollector(msg *discordgo.Message) {
	collectorsLock.Lock()
	collector, exists := Collectors[msg.ID]
	if exists {
		collector.deleteChan <- true
		delete(Collectors, msg.ID)
	}
	collectorsLock.Unlock()
}

func GetMessageComponentCollector(msg *discordgo.Message) (*Collector, bool) {
	collectorsLock.RLock()
	collector, exists := Collectors[msg.ID]
	collectorsLock.RUnlock()

	return collector, exists
}

func CleanComponentCollectors() {
	collectorsLock.Lock()
	defer collectorsLock.Unlock()

	for id, collector := range Collectors {
		if collector.Timeout != nil && time.Now().After(*collector.Timeout) {
			delete(Collectors, id)
		}
	}
}
