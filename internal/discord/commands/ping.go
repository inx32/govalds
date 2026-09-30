package commands

import (
	"context"
	"fmt"
	"time"

	"forge.pi.home.arpa/govalds/bot/internal/discord/gateway"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func ping(h *gateway.CommandHandler, ctx context.Context, event *events.ApplicationCommandInteractionCreate) {
	gatewayLatency := event.Client().Gateway.Latency()

	start := time.Now()
	event.DeferCreateMessage(false)
	elapsed := time.Since(start)

	response, err := event.Client().Rest.GetInteractionResponse(
		event.ApplicationID(), event.Token())
	if err != nil {
		panic("Unable to get interaction response: " + err.Error())
	}

	content := fmt.Sprintf(
		"## :man_in_manual_wheelchair: Линус GOвальдс работает!\n"+
			"`Gateway: %d мс.` `REST: %d мс.`\n"+
			"Защищаю людей от RedHat уже **%s**",
		gatewayLatency.Milliseconds(),
		elapsed.Milliseconds(),
		time.Since(h.Options().STime).Truncate(time.Second).String(),
	)

	event.Client().Rest.UpdateFollowupMessage(
		event.ApplicationID(), event.Token(), response.ID,
		discord.MessageUpdate{Content: &content},
	)
}

func Ping() gateway.Command {
	return gateway.Command{
		Name: "ping",
		Create: discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Send ping",
			IntegrationTypes: []discord.ApplicationIntegrationType{
				discord.ApplicationIntegrationTypeGuildInstall,
			},
		},
		Func: ping,
	}
}
