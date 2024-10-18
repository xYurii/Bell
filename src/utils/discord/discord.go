package discord

import (
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
)

type MessageBuilder struct {
	Session    *discordgo.Session
	ChannelID  string
	ReplyToID  string
	Content    string
	Embeds     []*discordgo.MessageEmbed
	Files      []*discordgo.File
	Components []discordgo.MessageComponent
	TTS        bool
}

func NewMessage(session *discordgo.Session, channelID, replyToID string) *MessageBuilder {
	return &MessageBuilder{
		Session:   session,
		ChannelID: channelID,
		ReplyToID: replyToID,
	}
}

func (m *MessageBuilder) WithButton(label string, customID string, style discordgo.ButtonStyle, emoji *discordgo.ComponentEmoji) *MessageBuilder {
	button := discordgo.Button{
		Label:    label,
		CustomID: customID,
		Style:    style,
		Emoji:    emoji,
	}
	if len(m.Components) > 0 {
		lastComponent := m.Components[len(m.Components)-1]

		if actionRow, ok := lastComponent.(discordgo.ActionsRow); ok {
			allButtons := true
			for _, component := range actionRow.Components {
				if _, isButton := component.(discordgo.Button); !isButton {
					allButtons = false
					break
				}
			}

			if allButtons && len(actionRow.Components) < 5 {
				actionRow.Components = append(actionRow.Components, button)
				m.Components[len(m.Components)-1] = actionRow
				return m
			}
		}
	}

	newActionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{button},
	}
	m.Components = append(m.Components, newActionRow)
	return m
}

func (m *MessageBuilder) WithSelectMenu(customID string, placeholder string, options []discordgo.SelectMenuOption, minValues int, maxValues int, disabled bool) *MessageBuilder {
	selectMenu := discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		Options:     options,
		MinValues:   &minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
	}
	m.Components = append(m.Components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{selectMenu},
	})
	return m
}

func (m *MessageBuilder) WithRoleSelect(customID string, placeholder string, minValues int, maxValues int, disabled bool) *MessageBuilder {
	selectMenu := discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		MenuType:    discordgo.RoleSelectMenu,
		MinValues:   &minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
	}
	m.Components = append(m.Components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{selectMenu},
	})
	return m
}

func (m *MessageBuilder) WithUserSelect(customID string, placeholder string, minValues int, maxValues int, disabled bool) *MessageBuilder {
	selectMenu := discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		MenuType:    discordgo.UserSelectMenu,
		MinValues:   &minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
	}
	m.Components = append(m.Components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{selectMenu},
	})
	return m
}

func (m *MessageBuilder) WithChannelSelect(customID string, placeholder string, minValues int, maxValues int, disabled bool) *MessageBuilder {
	selectMenu := discordgo.SelectMenu{
		CustomID:    customID,
		Placeholder: placeholder,
		MenuType:    discordgo.ChannelSelectMenu,
		MinValues:   &minValues,
		MaxValues:   maxValues,
		Disabled:    disabled,
	}
	m.Components = append(m.Components, discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{selectMenu},
	})
	return m
}

func (m *MessageBuilder) WithContent(content string) *MessageBuilder {
	m.Content = content
	return m
}

func (m *MessageBuilder) WithEmbed(embed *discordgo.MessageEmbed) *MessageBuilder {
	m.Embeds = append(m.Embeds, embed)
	return m
}

func (m *MessageBuilder) WithEmbeds(embeds []*discordgo.MessageEmbed) *MessageBuilder {
	m.Embeds = append(m.Embeds, embeds...)
	return m
}

func (m *MessageBuilder) WithFile(name string, reader io.Reader) *MessageBuilder {
	m.Files = append(m.Files, &discordgo.File{
		Name:        name,
		ContentType: "application/octet-stream",
		Reader:      reader,
	})
	return m
}

func (m *MessageBuilder) WithFiles(files []*discordgo.File) *MessageBuilder {
	m.Files = append(m.Files, files...)
	return m
}

func (m *MessageBuilder) WithTTS() *MessageBuilder {
	m.TTS = true
	return m
}

func (m *MessageBuilder) Send() (*discordgo.Message, error) {
	components := m.Components
	if len(components) > 5 {
		components = components[:5]
		log.Println("Too many components, only 5 are allowed!\nChannelID:", m.ChannelID, "\nMessageID:", m.ReplyToID)
	}
	msgSend := &discordgo.MessageSend{
		Content:    m.Content,
		Embeds:     m.Embeds,
		Files:      m.Files,
		TTS:        m.TTS,
		Components: components,
	}

	if m.ReplyToID != "" {
		msgSend.Reference = &discordgo.MessageReference{
			MessageID: m.ReplyToID,
			ChannelID: m.ChannelID,
		}
	}

	return m.Session.ChannelMessageSendComplex(m.ChannelID, msgSend)
}
