package commands

import "github.com/urfave/cli/v3"

var flagConfig = &cli.StringFlag{
	Name:     "config",
	Aliases:  []string{"c"},
	Usage:    "Path to config file",
	Required: true,
	OnlyOnce: true,
}

var flagEnv = &cli.StringSliceFlag{
	Name:    "env",
	Aliases: []string{"e"},
	Usage:   "Path to .env file",
}

var flagPid = &cli.StringFlag{
	Name:     "pid-file",
	Aliases:  []string{"p"},
	Usage:    "Path of a PID file",
	OnlyOnce: true,
}
