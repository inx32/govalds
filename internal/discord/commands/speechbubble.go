package commands

import (
	"bytes"
	"context"
	"errors"
	"image/png"

	"forge.pi.home.arpa/govalds/bot/internal/discord/gateway"
	"forge.pi.home.arpa/govalds/bot/internal/discord/util"
	"forge.pi.home.arpa/govalds/bot/pkg/imaging"
	"forge.pi.home.arpa/govalds/bot/pkg/imaging/bubble"
	"forge.pi.home.arpa/govalds/bot/pkg/imaging/bubbles"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func speechbubble(_ *gateway.CommandHandler, ctx context.Context, event *events.ApplicationCommandInteractionCreate) {
	event.DeferCreateMessage(false)
	data := event.SlashCommandInteractionData()

	optCorner := data.String("corner")
	optBubble := data.String("bubble")
	optExtend := data.Bool("extend")
	optHeight := data.Int("height")
	optNoGif := data.Bool("no_gif")

	if optCorner == "" {
		optCorner = "top_right"
	}
	if optBubble == "" {
		optBubble = "default"
	}
	if optHeight == 0 {
		optHeight = 20
	}

	corner := bubble.TopRight

	switch optCorner {
	case "top_left":
		corner = bubble.TopLeft

	case "bottom_right":
		corner = bubble.BottomRight

	case "bottom_left":
		corner = bubble.BottomLeft
	}

	imageNRGBA, err := util.DecodeImageAttachment(ctx, event, "image")
	if err != nil {
		var msg string

		switch {
		case errors.Is(err, imaging.ErrFileTooBig):
			msg = "Файл слишком большой (но не больше мамы цпловера)"

		case errors.Is(err, imaging.ErrImageTooBig):
			msg = "Изображение слишком большое (но не больше мамы цпловера)"

		case errors.Is(err, imaging.ErrImageInvalid):
			msg = "Не могу разобрать это изображение"

		default:
			panic(err)
		}

		util.FastError(event, msg)
		return
	}

	b := bubbles.Bubbles[optBubble]

	imageOutput, err := bubble.Speechbubble(imageNRGBA, b.Mask, corner, optExtend, optHeight)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	png.Encode(&buf, imageOutput)

	iresponse, err := event.Client().Rest.GetInteractionResponse(event.ApplicationID(), event.Token())
	if err != nil {
		panic(err)
	}

	filename := "image.gif"
	if optNoGif {
		filename = "image.png"
	}

	file := &discord.File{
		Name:   filename,
		Reader: &buf,
	}

	event.Client().Rest.UpdateFollowupMessage(
		event.ApplicationID(), event.Token(), iresponse.ID,
		discord.MessageUpdate{
			Files: []*discord.File{file},
		},
	)
}

func Speechbubble() gateway.Command {
	heightMin := 5
	heightMax := 100

	return gateway.Command{
		Name: "speechbubble",
		Create: discord.SlashCommandCreate{
			Name:        "speechbubble",
			Description: "Speechbubble",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionAttachment{
					Name:        "image",
					Description: "Source image",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name: "corner",
					NameLocalizations: map[discord.Locale]string{
						discord.LocaleRussian: "угол",
					},
					Description: "Corner where bubble will be added",
					Choices: []discord.ApplicationCommandOptionChoiceString{
						{
							Name:  "Top right",
							Value: "top_right",
						},
						{
							Name:  "Top left",
							Value: "top_left",
						},
						{
							Name:  "Bottom right",
							Value: "bottom_right",
						},
						{
							Name:  "Bottom left",
							Value: "bottom_left",
						},
					},
				},
				discord.ApplicationCommandOptionString{
					Name:        "bubble",
					Description: "Bubble asset",
					Choices:     bubbles.GetChoices(),
				},
				discord.ApplicationCommandOptionBool{
					Name:        "extend",
					Description: "Add an empty space for a bubble",
					Required:    false,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "height",
					Description: "Height of bubble (in % of source image height, default: 20%)",
					MinValue:    &heightMin,
					MaxValue:    &heightMax,
					Required:    false,
				},
				discord.ApplicationCommandOptionBool{
					Name:        "no_gif",
					Description: "Don't convert output image to GIF",
					Required:    false,
				},
			},
			IntegrationTypes: []discord.ApplicationIntegrationType{
				discord.ApplicationIntegrationTypeGuildInstall,
				discord.ApplicationIntegrationTypeUserInstall,
			},
		},
		Func: speechbubble,
	}
}
