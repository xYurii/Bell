package adapter

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type UserAdapter struct {
	Db *bun.DB
}

func NewUserAdapter(db *bun.DB) UserAdapter {
	return UserAdapter{Db: db}
}

func (a *UserAdapter) GetUser(ctx context.Context, id string, relations ...string) (user schemas.User) {
	query := a.Db.NewSelect().Model(&user).Where("id = ?", id)

	for _, relation := range relations {
		query.Relation(relation)
	}

	query.Scan(ctx)

	if user.ID == "" {
		user.ID = id
		a.CreateUser(ctx, user)
		user = a.GetUser(ctx, id, relations...)
	}

	return
}

func (a *UserAdapter) CreateUser(ctx context.Context, user schemas.User) error {
	_, err := a.Db.NewInsert().Model(&user).Exec(ctx)
	return err
}
