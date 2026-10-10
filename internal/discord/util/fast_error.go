package util

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func FastError(event *events.ApplicationCommandInteractionCreate, msg string) {
	iresponse, err := event.Client().Rest.GetInteractionResponse(event.ApplicationID(), event.Token())
	if err != nil {
		panic(err)
	}

	_, err = event.Client().Rest.UpdateFollowupMessage(
		event.ApplicationID(), event.Token(), iresponse.ID,
		discord.MessageUpdate{Content: &msg},
	)
	if err != nil {
		panic(err)
	}
}
