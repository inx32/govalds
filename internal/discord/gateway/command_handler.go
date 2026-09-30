package gateway

import (
	"context"
	"errors"
	"fmt"
	"time"

	"forge.pi.home.arpa/govalds/bot/pkg/concurrency"
	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/rs/zerolog"
)

func logconfig(logevent *zerolog.Event, event *events.ApplicationCommandInteractionCreate, eventID string) *zerolog.Event {
	return logevent.
		Str("command_name", event.Data.CommandName()).
		Str("user_id", event.User().ID.String()).
		Str("channel_id", event.Channel().ID().String()).
		Bool("is_guild", event.GuildID() != nil).
		Str("event_id", eventID)
}

type CommandFactory = func() Command

type Command struct {
	Name   string
	Create discord.SlashCommandCreate
	Func   func(*CommandHandler, context.Context, *events.ApplicationCommandInteractionCreate)
}

type CommandHandlerOptions struct {
	STime   time.Time
	Context context.Context
	Config  *config.Config
	Logger  zerolog.Logger
}

type CommandHandler struct {
	opts     CommandHandlerOptions
	commands *concurrency.Map[string, Command]
}

func (h *CommandHandler) Options() CommandHandlerOptions { return h.opts }

func (h *CommandHandler) RegisterCommand(cmd Command) error {
	if h.commands.Has(cmd.Name) {
		return errors.New("this command is already registered")
	}
	h.commands.Set(cmd.Name, cmd)
	h.opts.Logger.Trace().
		Str("name", cmd.Name).
		Msg("Registered application command")
	return nil
}

func (h *CommandHandler) Sync(client *bot.Client) error {
	var creates []discord.ApplicationCommandCreate
	for _, cmd := range h.commands.Iter() {
		creates = append(creates, cmd.Create)
	}
	_, err := client.Rest.SetGlobalCommands(client.ApplicationID, creates)
	if err == nil {
		h.opts.Logger.Debug().Msg("Synced global application commands")
	}
	return err
}

func (h *CommandHandler) Process(ctx context.Context, event *events.ApplicationCommandInteractionCreate) {
	handler, ok := h.commands.Get(event.Data.CommandName())
	if !ok || handler.Func == nil {
		eventID := ctx.Value("event_id").(string)

		logconfig(h.opts.Logger.Error(), event, eventID).
			Msg("Handler for application command is not found")

		event.CreateMessage(discord.MessageCreate{
			Content: fmt.Sprintf(
				":warning: __**Something went wrong!**__\n"+
					"```Handler for this application command not found```\n"+
					"This incident has been reported. Event ID: `%s`",
				eventID,
			),
		})
		return
	}

	defer func() {
		if err2 := recover(); err2 != nil {
			eventID := ctx.Value("event_id").(string)

			event.DeferCreateMessage(false)
			response, err := event.Client().Rest.GetInteractionResponse(
				event.ApplicationID(), event.Token())
			if err != nil {
				logconfig(h.opts.Logger.Err(err), event, eventID).
					Msg("Handle panic: unable to get interaction response")
				panic(err2)
			}

			content := fmt.Sprintf(
				":warning: __**Something went wrong!**__\n"+
					"```Application command handler panicked!```\n"+
					"This incident has been reported. Event ID: `%s`",
				eventID,
			)

			event.Client().Rest.UpdateFollowupMessage(
				event.ApplicationID(), event.Token(), response.ID,
				discord.MessageUpdate{Content: &content},
			)

			panic(err2)
		}
	}()

	handler.Func(h, ctx, event)
}

func NewCommandHandler(opts CommandHandlerOptions) *CommandHandler {
	handler := &CommandHandler{
		opts:     opts,
		commands: concurrency.NewMap[string, Command](),
	}
	return handler
}
