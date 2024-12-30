package discord

import "github.com/bwmarrin/discordgo"

type Button struct {
	discordgo.Button
}

func NewButton() *Button {
	return &Button{
		Button: discordgo.Button{},
	}
}

func (b *Button) WithLabel(label string) *Button {
	b.Button.Label = label
	return b
}

func (b *Button) WithCustomID(customID string) *Button {
	b.Button.CustomID = customID
	return b
}

func (b *Button) WithStyle(style discordgo.ButtonStyle) *Button {
	b.Button.Style = style
	return b
}

func (b *Button) WithEmoji(emoji *discordgo.ComponentEmoji) *Button {
	b.Button.Emoji = emoji
	return b
}

func (b *Button) WithURL(url string) {
	b.Button.URL = url
}

func (b *Button) Build() *discordgo.Button {
	return &b.Button
}
