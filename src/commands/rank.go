package commands

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"github.com/xYurii/Bell/src/database"
	"github.com/xYurii/Bell/src/database/schemas"
	"github.com/xYurii/Bell/src/handler"
	"github.com/xYurii/Bell/src/utils"
	"github.com/xYurii/Bell/src/utils/discord"
)

const (
	width      = 400
	height     = 510
	avatarSize = 80
	fontSize   = 18
	background = "#FCE07D"
)

func init() {
	handler.RegisterCommand(handler.Command{
		Name:      "top",
		Aliases:   []string{"rank"},
		Run:       runTop,
		Developer: true,
	})
}

func runTop(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	users := database.User.SortUsers(ctx, 5, "status_time")
	img, err := generateImage(s, users)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	discord.NewMessage(s, m.ChannelID, m.ID).WithFile("top.png", &buf).Send()
}

func generateImage(s *discordgo.Session, users []*schemas.User) (image.Image, error) {
	dc := gg.NewContext(width, height)
	dc.SetHexColor(background)
	dc.Clear()

	if err := dc.LoadFontFace("NexaBold.ttf", fontSize); err != nil {
		return nil, err
	}

	y := 20
	for _, data := range users {
		user, err := s.User(data.ID)
		if user == nil {
			return nil, err
		}

		avatar, err := downloadAndResizeAvatar(user.AvatarURL("2048"), avatarSize)
		if err != nil {
			log.Println("Erro ao baixar avatar:", err)
			continue
		}

		dc.DrawCircle(50, float64(y+avatarSize/2), avatarSize/2)
		dc.Clip()
		dc.DrawImage(avatar, 10, y)
		dc.ResetClip()

		dc.SetRGB(0, 0, 0)
		dc.DrawString(user.Username, 100, float64(y+20))
		dc.DrawString(fmt.Sprintf("ID: %s", user.ID), 100, float64(y+40))
		dc.DrawString(fmt.Sprintf("Tempo: %s", utils.FormatDuration(data.StatusTime)), 100, float64(y+60))

		y += avatarSize + 20
	}
	return dc.Image(), nil
}

func downloadAndResizeAvatar(url string, size int) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return resize.Resize(uint(size), uint(size), img, resize.Lanczos3), nil
}
