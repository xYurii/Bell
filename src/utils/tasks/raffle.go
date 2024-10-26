package tasks

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type RaffleManager struct {
	Db      *bun.DB
	Session *discordgo.Session
	lock    sync.Mutex
}

func NewRaffleManager(db *bun.DB, s *discordgo.Session) *RaffleManager {
	return &RaffleManager{
		Db:      db,
		Session: s,
	}
}

func (rm *RaffleManager) EndExpiredRaffles() {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	var raffles []schemas.Raffle
	now := time.Now()
	ctx := context.Background()

	err := rm.Db.NewSelect().
		Model(&raffles).
		Where("ends_at < ?", now).
		Where("ended_at IS NULL").
		Order("ends_at DESC").
		Scan(ctx)

	if err != nil {
		log.Println("Error while fetching raffles:", err)
		return
	}

	for _, raffle := range raffles {
		log.Println("Ending raffle:", raffle.ID, "(", raffle.RaffleType.String(), ")")
		err := rm.EndRaffle(ctx, &raffle)
		if err != nil {
			log.Fatalln("Error ending raffle:", err)
			continue
		}

		rm.CreateNewRaffle(ctx, raffle)
	}
}

func (rm *RaffleManager) EndRaffle(ctx context.Context, raffle *schemas.Raffle) error {
	tx, err := rm.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	totalTickets, err := tx.NewSelect().
		Model((*schemas.RaffleTickets)(nil)).
		Where("id = ?", raffle.ID).
		Count(ctx)
	if err != nil {
		return fmt.Errorf("error counting tickets: %w", err)
	}

	fmt.Println("Total tickets:", totalTickets)

	if totalTickets == 0 {
		newEnd := raffle.RaffleType.EndTime()
		raffle.EndsAt = newEnd
		_, err = tx.NewUpdate().
			Model(raffle).
			Where("id = ?", raffle.ID).
			Column("ends_at").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("error updating ends_at: %w", err)
		}
		return tx.Commit()
	}

	winnerTicket, err := rm.SelectWinnerTicket(ctx, &tx, raffle.ID, totalTickets)
	if err != nil {
		return err
	}

	fmt.Println("Winner:", winnerTicket.UserID, "Total tickets:", totalTickets)

	now := time.Now()
	winnerTicketID := winnerTicket.ID
	raffle.WinnerTicketID = winnerTicketID
	raffle.EndedAt = now

	_, err = tx.NewUpdate().
		Model(raffle).
		Where("id = ?", raffle.ID).
		Column("ended_at", "winner_ticket_id").
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error ending raffle: %w", err)
	}

	return tx.Commit()
}

func (rm *RaffleManager) SelectWinnerTicket(ctx context.Context, tx *bun.Tx, raffleID, totalTickets int) (*schemas.RaffleTickets, error) {
	skip := rand.Intn(totalTickets)
	var winnerTicket schemas.RaffleTickets

	err := tx.NewSelect().
		Model(&winnerTicket).
		Where("raffle_id = ?", raffleID).
		Order("bought_at ASC").
		Offset(skip).
		Limit(1).
		For("UPDATE SKIP LOCKED").
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error selecting winner: %w", err)
	}

	return &winnerTicket, nil
}

func (rm *RaffleManager) CreateNewRaffle(ctx context.Context, raffle schemas.Raffle) error {
	now := time.Now()
	newRaffle := &schemas.Raffle{
		RaffleType: raffle.RaffleType,
		StartedAt:  now,
		EndsAt:     raffle.RaffleType.EndTime(),
	}

	_, err := rm.Db.NewInsert().Model(newRaffle).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating new raffle: %w", err)
	}

	return nil
}

func (rm *RaffleManager) EnsureRafflesExist(ctx context.Context) error {
	tx, err := rm.Db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	raffleTypes := []schemas.RaffleType{
		schemas.DailyRaffle,
		schemas.NormalRaffle,
		schemas.Lightning,
	}

	for _, raffleType := range raffleTypes {
		count, err := tx.NewSelect().
			Model((*schemas.Raffle)(nil)).
			Where("raffle_type = ?", raffleType).
			Count(ctx)
		if err != nil {
			return fmt.Errorf("error checking raffle type %s: %w", raffleType.String(), err)
		}

		if count == 0 {
			newRaffle := &schemas.Raffle{
				RaffleType: raffleType,
				StartedAt:  time.Now(),
				EndsAt:     raffleType.EndTime(),
			}

			_, err = tx.NewInsert().Model(newRaffle).Exec(ctx)
			if err != nil {
				return fmt.Errorf("error creating raffle for type %s: %w", raffleType.String(), err)
			}
			log.Printf("Created a new raffle for type %s", raffleType.String())
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
