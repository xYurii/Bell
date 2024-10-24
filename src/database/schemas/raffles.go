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

func (rt RaffleType) Price() int {
	return [3]int{250, 500, 1000}[rt]
}

func (rt RaffleType) MaxTickets() int {
	base := 100000
	max := base * int(rt+1)
	return max
}

type Raffle struct {
	bun.BaseModel `bun:"table:raffles,alias:r"`

	ID             int        `bun:"id,pk,autoincrement"`
	RaffleType     RaffleType `bun:"raffle_type,notnull"`
	StartedAt      time.Time  `bun:"started_at,notnull,default:current_timestamp"`
	EndsAt         time.Time  `bun:"ends_at,notnull"`
	EndedAt        time.Time  `bun:"ended_at,nullzero"`
	WinnerTicketID *int64     `bun:"winner_ticket_id,nullzero"`
}

type RaffleTickets struct {
	bun.BaseModel `bun:"table:raffle_tickets,alias:rt"`

	ID       int       `bun:"id,pk,autoincrement"`
	BoughtAt time.Time `bun:"bought_at,notnull,default:current_timestamp"`
	UserID   string    `bun:"user_id,notnull"`
}
