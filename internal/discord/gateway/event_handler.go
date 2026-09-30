package gateway

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/rs/zerolog"
)

type EventHandlerOptions struct {
	STime   time.Time
	Context context.Context
	Config  *config.Config
	Logger  zerolog.Logger
	Tasks   *sync.WaitGroup
}

type EventHandler struct {
	opts EventHandlerOptions

	startupTime    time.Time
	commandHandler *CommandHandler
	onReadyOnce    sync.Once
}

func (h *EventHandler) OnEvent(event bot.Event) {
	seq := event.Client().Gateway.LastSequenceReceived()
	id := strconv.FormatUint(uint64(h.startupTime.UnixNano()), 10) +
		":" + strconv.FormatInt(int64(*seq), 10)

	ctx := context.WithValue(h.opts.Context, "event_id", id)

	defer func() {
		if err := recover(); err != nil {
			h.opts.Logger.Error().
				Str("event_id", id).
				Str("event_type", fmt.Sprintf("%T", event)).
				Any("err", err).
				Msgf("Got panic in event handler! Stack:\n%s", string(debug.Stack()))
		}
	}()

	switch e := event.(type) {
	case *events.Ready:
		h.onReady(ctx, e)

	case *events.ApplicationCommandInteractionCreate:
		h.commandHandler.Process(ctx, e)
	}
}

func NewEventHandler(opts EventHandlerOptions, commands []CommandFactory) (*EventHandler, error) {
	cmd := NewCommandHandler(CommandHandlerOptions{
		STime:   opts.STime,
		Context: opts.Context,
		Config:  opts.Config,
		Logger:  opts.Logger,
	})

	for _, cfunc := range commands {
		c := cfunc()
		if err := cmd.RegisterCommand(c); err != nil {
			return nil, fmt.Errorf("register command %s: %w", c.Name, err)
		}
	}

	return &EventHandler{
		opts:           opts,
		startupTime:    time.Now(),
		commandHandler: cmd,
	}, nil
}
