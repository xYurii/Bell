package schemas

import (
	"time"

	"github.com/uptrace/bun"
)

type Guild struct {
	bun.BaseModel `bun:"table:guilds,alias:g"`

	ID        string    `bun:"id,pk,notnull,unique"`
	Prefix    string    `bun:"prefix,default:.."`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
}
