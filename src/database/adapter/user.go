package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type UserAdapter struct {
	Db *bun.DB
}

func NewUserAdapter(db *bun.DB) UserAdapter {
	return UserAdapter{Db: db}
}

func (a *UserAdapter) GetUser(ctx context.Context, author *discordgo.User, relations ...string) (user schemas.User) {
	query := a.Db.NewSelect().Model(&user).Where("id = ?", author.ID)

	for _, relation := range relations {
		query.Relation(relation)
	}

	query.Scan(ctx)

	if user.ID == "" {
		user.ID = author.ID
		a.CreateUser(ctx, user)
		user = a.GetUser(ctx, author, relations...)
	}

	return
}

func (a *UserAdapter) CreateUser(ctx context.Context, user schemas.User) error {
	_, err := a.Db.NewInsert().Model(&user).Exec(ctx)
	return err
}

func (a *UserAdapter) UpdateUser(ctx context.Context, author *discordgo.User, callback func(user schemas.User) schemas.User, relations ...string) error {
	return a.Db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tx.ExecContext(ctx, fmt.Sprintf("SELECT pg_advisory_xact_lock(%s)", author.ID))
		user := a.GetUser(ctx, author, relations...)
		user = callback(user)
		_, err := a.Db.NewUpdate().Model(&user).Where("id = ?", author.ID).Exec(ctx)
		return err
	})
}
