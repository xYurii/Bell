package discord

import (
	"github.com/bwmarrin/discordgo"
)

func MemberHasPermission(s *discordgo.Session, m *discordgo.Message, permission int64) (bool, error) {
	permissions, err := s.State.MessagePermissions(m)
	if err != nil {
		return false, err
	}

	hasPerm := permissions&permission != 0
	return hasPerm, nil
}
