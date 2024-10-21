package discord

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	Embed *discordgo.MessageEmbed
}

func NewEmbedBuilder() *EmbedBuilder {
	return &EmbedBuilder{
		Embed: &discordgo.MessageEmbed{},
	}
}

func (e *EmbedBuilder) WithURL(url string) *EmbedBuilder {
	e.Embed.URL = url
	return e
}

func (e *EmbedBuilder) WithAuthor(name string, iconURL string, url string) *EmbedBuilder {
	e.Embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
		URL:     url,
	}
	return e
}

func (e *EmbedBuilder) WithFields(fields ...*discordgo.MessageEmbedField) *EmbedBuilder {
	e.Embed.Fields = append(e.Embed.Fields, fields...)
	return e
}

func (e *EmbedBuilder) WithTimestamp(timestamp string) *EmbedBuilder {
	e.Embed.Timestamp = timestamp
	return e
}

func (e *EmbedBuilder) WithField(name string, value string, inline bool) *EmbedBuilder {
	e.Embed.Fields = append(e.Embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return e
}

func (e *EmbedBuilder) WithTitle(title string) *EmbedBuilder {
	e.Embed.Title = title
	return e
}

func (e *EmbedBuilder) WithDescription(description string) *EmbedBuilder {
	e.Embed.Description = description
	return e
}

func (e *EmbedBuilder) WithColor(color int) *EmbedBuilder {
	e.Embed.Color = color
	return e
}

func (e *EmbedBuilder) WithFooter(text string, iconURL string) *EmbedBuilder {
	e.Embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return e
}

func (e *EmbedBuilder) WithImage(url string) *EmbedBuilder {
	e.Embed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) WithThumbnail(url string) *EmbedBuilder {
	e.Embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) Build() *discordgo.MessageEmbed {
	if len(e.Embed.Fields) > 25 {
		e.Embed.Fields = e.Embed.Fields[:25]
	}
	return e.Embed
}
