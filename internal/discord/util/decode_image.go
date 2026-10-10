package util

import (
	"context"
	"image"
	"net/http"
	"time"

	"forge.pi.home.arpa/govalds/bot/pkg/imaging"
	"github.com/disgoorg/disgo/events"
)

func DecodeImageAttachment(ctx context.Context, event *events.ApplicationCommandInteractionCreate, name string) (*image.NRGBA, error) {
	data := event.SlashCommandInteractionData()
	optImage := data.Attachment(name)

	requestCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, "GET", optImage.URL, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	img, err := imaging.Decode(response.Body)
	if err != nil {
		response.Body.Close()
		return nil, err
	}
	response.Body.Close()

	return imaging.ToNRGBA(img), nil
}
