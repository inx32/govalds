package commands

import (
	"context"
	"fmt"
	"time"

	"forge.pi.home.arpa/govalds/bot/internal/discord/gateway"
	"forge.pi.home.arpa/govalds/bot/internal/versioninfo"
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
		"## :man_in_manual_wheelchair: Я работаю!\n"+
			"Линус GOвальдс версия %s, коммит `%s`, собран <t:%d:f>\n"+
			"Репозиторий: `%s`\n\n"+
			"`Gateway: %d мс.` `REST: %d мс.`\n"+
			"Защищаю людей от RedHat уже **%s**",
		versioninfo.BotVersion(),
		versioninfo.Commit(),
		versioninfo.BuildTime().Unix(),
		versioninfo.Repository(),
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
				discord.ApplicationIntegrationTypeUserInstall,
			},
		},
		Func: ping,
	}
}
