package schemas

import (
	"time"

	"github.com/uptrace/bun"
)

type CommandsUsed struct {
	bun.BaseModel `bun:"table:commands_used,alias:cu"`

	ID             int       `bun:"id,pk,autoincrement"`
	UserID         string    `bun:"user_id,notnull"`
	GuildID        string    `bun:"guild_id,notnull"`
	ChannelID      string    `bun:"channel_id,notnull"`
	CommandName    string    `bun:"command_name,notnull"`
	MessageContent string    `bun:"message_content,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp"`
}
