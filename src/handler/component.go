package handler

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Component struct {
	Name string
	Type discordgo.ComponentType
	Run  func(context.Context, *discordgo.Session, *discordgo.InteractionCreate)
}

var components = map[string]Component{}

func RegisterComponent(component Component) {
	components[component.Name] = component
}

func GetComponent(customId string) (Component, bool) {
	customIdParsed := ParseComponentId(customId)
	component, exists := components[customIdParsed[0]]
	return component, exists
}

func ParseComponentId(id string) []string {
	split := strings.Split(id, "_")
	return split
}
