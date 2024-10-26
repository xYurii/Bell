package schemas

import (
	"time"

	"github.com/uptrace/bun"
)

type RaffleType int

const (
	DailyRaffle RaffleType = iota
	NormalRaffle
	Lightning
)

func (rt RaffleType) String() string {
	return [3]string{"DAILY", "NORMAL", "LIGHTNING"}[rt]
}

func (rt RaffleType) Price() int {
	return [3]int{250, 500, 1000}[rt]
}

func (rt RaffleType) MaxTickets() int {
	base := 100000
	max := base * int(rt+1)
	return max
}

func (rt RaffleType) EndTime() time.Time {
	location, _ := time.LoadLocation("America/Sao_Paulo")
	now := time.Now().In(location)

	switch rt {
	case DailyRaffle:
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, location)
		return nextMidnight
	case NormalRaffle:
		return now.Add(1 * time.Hour)
	case Lightning:
		return now.Add(15 * time.Minute)
	default:
		return now
	}
}

type Raffle struct {
	bun.BaseModel `bun:"table:raffles,alias:r"`

	ID             int            `bun:"id,pk,autoincrement"`
	RaffleType     RaffleType     `bun:"raffle_type,notnull"`
	StartedAt      time.Time      `bun:"started_at,notnull,default:current_timestamp"`
	EndsAt         time.Time      `bun:"ends_at,notnull"`
	EndedAt        time.Time      `bun:"ended_at,nullzero"`
	WinnerTicket   *RaffleTickets `bun:"rel:belongs-to,join:winner_ticket_id=id"`
	WinnerTicketID int            `bun:"winner_ticket_id,nullzero"`
	RewardPrice    int64          `bun:"reward_price,notnull,default:0"`
}

type RaffleTickets struct {
	bun.BaseModel `bun:"table:raffle_tickets,alias:rt"`

	ID       int       `bun:"id,pk,autoincrement"`
	BoughtAt time.Time `bun:"bought_at,notnull,default:current_timestamp"`
	UserID   string    `bun:"user_id,notnull"`
	RaffleID int       `bun:"raffle_id,notnull"`
}
