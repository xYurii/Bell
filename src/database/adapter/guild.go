package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type GuildAdapter struct {
	Db *bun.DB
}

func NewGuildAdapter(db *bun.DB) GuildAdapter {
	return GuildAdapter{Db: db}
}

func (a *GuildAdapter) GetGuild(ctx context.Context, id string, relations ...string) (guild schemas.Guild) {
	query := a.Db.NewSelect().Model(&guild).Where("id = ?", id)

	for _, relation := range relations {
		query.Relation(relation)
	}

	query.Scan(ctx)

	if guild.ID == "" {
		guild.ID = id
		guild.Prefix = ".."
		a.CreateGuild(ctx, guild)
		guild = a.GetGuild(ctx, id, relations...)
	}

	return
}

func (a *GuildAdapter) CreateGuild(ctx context.Context, guild schemas.Guild) error {
	_, err := a.Db.NewInsert().Model(&guild).Exec(ctx)
	return err
}

func (a *GuildAdapter) UpdateGuild(ctx context.Context, id string, callback func(guild schemas.Guild) schemas.Guild, relations ...string) error {
	return a.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%s)", id))
		g := a.GetGuild(ctx, id, relations...)
		g = callback(g)
		_, err := a.Db.NewUpdate().Model(&g).Where("id = ?", id).Exec(ctx)
		return err
	})
}
