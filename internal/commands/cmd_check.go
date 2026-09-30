package commands

import (
	"context"
	"fmt"

	"forge.pi.home.arpa/govalds/bot/pkg/config"
	"forge.pi.home.arpa/govalds/bot/pkg/dotenv"
	"github.com/urfave/cli/v3"
)

var Check = &cli.Command{
	Name:      "check",
	Usage:     "Check config",
	UsageText: "linus-bot check [options]",
	Action:    check,
	Flags:     []cli.Flag{flagConfig, flagEnv},
}

func check(ctx context.Context, command *cli.Command) error {
	if files := command.StringSlice("env"); len(files) != 0 {
		if err := dotenv.LoadFiles(files); err != nil {
			return fmt.Errorf("load env: %w", err)
		}
	}

	loc := command.String("config")
	fmt.Println("\033[1;97mOpening", loc, "\033[0m")
	_, err := config.Load(loc)
	if err != nil {
		return fmt.Errorf("\033[1;91m%w\033[0m", err)
	}
	fmt.Println("\033[1;92mConfig is valid, no error reported\033[0m")
	return nil
}
