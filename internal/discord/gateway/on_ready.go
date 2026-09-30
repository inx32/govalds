package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func onReadyPresenceLoop(ctx context.Context, client *bot.Client) {
	presences := []string{
		"новый проект redhat systemd вирус",
		"неопознанный летающий объект НЛО redhat",
		"go run ./cmd start -c config.yaml",
		"Слава Linux !!!",
	}

	var index int

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(120 * time.Second):
			presence := presences[index]
			client.SetPresence(ctx, gateway.WithCustomActivity(presence))

			index++
			if index >= len(presences) {
				index = 0
			}
		}
	}
}

func (h *EventHandler) onReady(ctx context.Context, event *events.Ready) {
	h.opts.Logger.Info().
		Str("user", fmt.Sprintf("%s#%s", event.User.Username, event.User.Discriminator)).
		Msg("Connected to gateway!")

	h.onReadyOnce.Do(func() {
		h.opts.Logger.Debug().Msg("Starting presence updater...")
		h.opts.Tasks.Go(func() {
			onReadyPresenceLoop(ctx, event.Client())
		})

		if err := h.commandHandler.Sync(event.Client()); err != nil {
			h.opts.Logger.Err(err).Msg("Unable to sync application commands")
		} else {
			h.opts.Logger.Info().Msg("Synced application commands")
		}
	})
}
