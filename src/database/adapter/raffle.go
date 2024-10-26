package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type RaffleAdapter struct {
	Db *bun.DB
}

type ActualRaffleStatus struct {
	Raffle         schemas.Raffle
	PreviousRaffle schemas.Raffle
	TicketsCount   int
	UsersCount     int
	ActualReward   int64
}

func NewRaffleAdapter(db *bun.DB) RaffleAdapter {
	return RaffleAdapter{Db: db}
}

func (a *RaffleAdapter) GetRaffleUsersCount(ctx context.Context, raffleID int, raffleType schemas.RaffleType) (raffleUsers []schemas.RaffleTickets, err error) {
	err = a.Db.NewSelect().
		Model(&raffleUsers).
		ColumnExpr("DISTINCT user_id").
		Where("raffle_id = ?", raffleID).
		Scan(ctx)

	if err != nil {
		return []schemas.RaffleTickets{}, err
	}

	return raffleUsers, nil
}

func (a *RaffleAdapter) GetRaffleStatus(ctx context.Context, raffleType schemas.RaffleType) (ActualRaffleStatus, error) {
	var previousRaffle schemas.Raffle
	var actualRaffle schemas.Raffle
	err := a.Db.NewSelect().
		Model(&actualRaffle).
		Where("raffle_type = ?", raffleType).
		Where("ended_at IS NULL").
		Limit(1).
		Scan(ctx)

	if err != nil {
		return ActualRaffleStatus{}, fmt.Errorf("error getting raffle status: %w", err)
	}

	err = a.Db.NewSelect().
		Model(&previousRaffle).
		Where("ended_at IS NOT NULL").
		Where("raffle_type = ?", raffleType).
		Order("ended_at DESC").
		Relation("WinnerTicket").
		Limit(1).
		Scan(ctx)

	if err != nil {
		previousRaffle = schemas.Raffle{
			RewardPrice: 0,
		}
	}

	ticketsCount, err := a.GetRaffleTicketsCount(ctx, actualRaffle.ID)
	if err != nil {
		return ActualRaffleStatus{}, fmt.Errorf("error getting tickets count: %w", err)
	}

	usersCount, err := a.GetRaffleUsersCount(ctx, actualRaffle.ID, raffleType)
	if err != nil {
		return ActualRaffleStatus{}, fmt.Errorf("error getting users count: %w", err)
	}

	actualReward := int64(len(usersCount)) * int64(raffleType.Price())

	return ActualRaffleStatus{
		Raffle:         actualRaffle,
		PreviousRaffle: previousRaffle,
		TicketsCount:   ticketsCount,
		UsersCount:     len(usersCount),
		ActualReward:   actualReward,
	}, err
}

func (a *RaffleAdapter) InsertTickets(ctx context.Context, raffleID int, userID string, ticketsCount int64, tx bun.Tx) error {
	tickets := make([]schemas.RaffleTickets, ticketsCount)
	date := time.Now()
	for i := int64(0); i < ticketsCount; i++ {
		tickets[i] = schemas.RaffleTickets{
			BoughtAt: date,
			UserID:   userID,
			RaffleID: raffleID,
		}
	}
	_, err := tx.NewInsert().Model(&tickets).Exec(ctx)
	return err
}

func (a *RaffleAdapter) GetRaffleTicketsCount(ctx context.Context, raffleID int) (int, error) {
	count, err := a.Db.NewSelect().
		Model((*schemas.RaffleTickets)(nil)).
		Where("raffle_id = ?", raffleID).
		Count(ctx)
	return count, err
}

// func (a *RaffleAdapter)
