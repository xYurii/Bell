package adapter

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/xYurii/Bell/src/database/schemas"
)

type CommandsUsedAdapter struct {
	Db *bun.DB
}

func NewCommandsUsedAdapter(db *bun.DB) CommandsUsedAdapter {
	return CommandsUsedAdapter{Db: db}
}

func (a *CommandsUsedAdapter) InsertCommand(ctx context.Context, commandName, userID, guildID, channelID, messageContent string) error {
	_, err := a.Db.NewInsert().
		Model(&schemas.CommandsUsed{
			CommandName:    commandName,
			UserID:         userID,
			GuildID:        guildID,
			ChannelID:      channelID,
			MessageContent: messageContent,
		}).
		Exec(ctx)

	return err
}
