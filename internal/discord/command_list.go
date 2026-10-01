package discord

import (
	"forge.pi.home.arpa/govalds/bot/internal/discord/commands"
	"forge.pi.home.arpa/govalds/bot/internal/discord/gateway"
)

var CommandList = []gateway.CommandFactory{commands.Ping, commands.Speechbubble}
