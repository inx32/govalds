package bot

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	idiscord "forge.pi.home.arpa/govalds/bot/internal/discord"
	igateway "forge.pi.home.arpa/govalds/bot/internal/discord/gateway"
	"forge.pi.home.arpa/govalds/bot/internal/versioninfo"
	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"forge.pi.home.arpa/govalds/bot/pkg/logging"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/rs/zerolog"
)

type linus struct {
	stime   time.Time
	config  *config.Config
	logger  zerolog.Logger
	tasks   *sync.WaitGroup
	client  *bot.Client
	wrapper *gatewayWrapper
}

func (l *linus) Start(ctx context.Context) error {
	l.logger.Info().
		Str("version", versioninfo.BotVersion()).
		Str("commit", versioninfo.Commit()).
		Msg("Starting Linus bot...")

	handler, err := igateway.NewEventHandler(igateway.EventHandlerOptions{
		STime:   l.stime,
		Context: ctx,
		Config:  l.config,
		Logger:  l.logger,
		Tasks:   l.tasks,
	}, idiscord.CommandList)
	if err != nil {
		return err
	}
	l.wrapper.next = handler

	if err := l.client.OpenGateway(ctx); err != nil {
		return fmt.Errorf("open gateway: %w", err)
	}

	<-ctx.Done()
	l.logger.Info().Msg("Context canceled, closing...")
	l.wrapper.Close()
	l.client.Close(context.TODO())

	l.logger.Info().Msg("Linus GOvalds closed successfully")
	return nil
}

func New(conf *config.Config, goModule string) (*linus, error) {
	logger, err := logging.NewLogger(conf.Logging, goModule)
	if err != nil {
		return nil, fmt.Errorf("unable to create logger: %w", err)
	}

	wrapper := &gatewayWrapper{log: *logger}

	clientLogger := logger.Level(zerolog.InfoLevel)
	if conf.Logging != nil && conf.Logging.ClientTrace {
		clientLogger = logger.Level(zerolog.TraceLevel)
	}
	clientLogger = clientLogger.With().Str("from", "client").Logger()

	client, err := disgo.New(conf.Token.Value(),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsNone),
			gateway.WithBrowser("Discord Android"),
		),
		bot.WithLogger(slog.New(zerolog.NewSlogHandler(clientLogger))),
		bot.WithEventListeners(wrapper),
	)

	if err != nil {
		return nil, err
	}

	return &linus{
		stime:   time.Now(),
		config:  conf,
		logger:  *logger,
		tasks:   new(sync.WaitGroup),
		client:  client,
		wrapper: wrapper,
	}, nil
}
