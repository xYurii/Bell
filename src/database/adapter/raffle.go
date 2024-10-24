package adapter

import "github.com/uptrace/bun"

type RaffleAdapter struct {
	Db *bun.DB
}

func NewRaffleAdapter(db *bun.DB) RaffleAdapter {
	return RaffleAdapter{Db: db}
}
