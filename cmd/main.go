package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"forge.pi.home.arpa/govalds/bot/internal/commands"
	"github.com/urfave/cli/v3"
)

var command = &cli.Command{
	Name:      "linus",
	Usage:     "Linus GOvalds bot",
	UsageText: "linus [command] [options]",
	Commands:  []*cli.Command{commands.Start, commands.Check},
}

func main() {
	cli.HelpFlag = nil

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	err := command.Run(ctx, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error: %s\n", command.Name, err)
		os.Exit(1)
	}
}
