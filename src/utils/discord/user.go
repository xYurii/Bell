package discord

import "github.com/bwmarrin/discordgo"

func GetUser(s *discordgo.Session, evt *discordgo.MessageCreate, args []string) *discordgo.User {
	if len(args) < 1 {
		return evt.Author
	}

	if len(evt.Mentions) > 0 {
		return evt.Mentions[0]
	}

	member, _ := s.State.Member(evt.GuildID, args[0])
	if member != nil {
		return member.User
	}

	user, err := s.User(args[0])
	if err != nil {
		return evt.Author
	}

	return user
}
