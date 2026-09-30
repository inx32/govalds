package commands

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"runtime/debug"
	"syscall"

	"forge.pi.home.arpa/govalds/bot/internal/bot"
	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"forge.pi.home.arpa/govalds/bot/pkg/dotenv"
	"forge.pi.home.arpa/govalds/bot/pkg/pidfile"
	"github.com/urfave/cli/v3"
)

var Start = &cli.Command{
	Name:      "start",
	Usage:     "Start bot",
	UsageText: "linus-bot start [options]",
	Action:    start,
	Flags:     []cli.Flag{flagConfig, flagEnv, flagPid},
}

func start(ctx context.Context, command *cli.Command) error {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("unable to read build info")
	}

	if pidFilename := command.String("pid-file"); pidFilename != "" {
		pid := pidfile.New(pidFilename)
		if err := pid.Lock(); err != nil {
			return fmt.Errorf("create pidfile: %w", err)
		}
		defer pid.Release()
	}

	if files := command.StringSlice("env"); len(files) != 0 {
		if err := dotenv.LoadFiles(files); err != nil {
			return fmt.Errorf("load env: %w", err)
		}
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conf, err := config.Load(command.String("config"))
	if err != nil {
		return err
	}

	linus, err := bot.New(conf, bi.Main.Path)
	if err != nil {
		return err
	}
	return linus.Start(ctx)
}
