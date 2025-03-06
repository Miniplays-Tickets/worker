package utils

import (
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/errorcontext"
	"github.com/rxdn/gdl/gateway/payloads/events"
)

func MessageCreateErrorContext(e events.MessageCreate) errorcontext.WorkerErrorContext {
	return errorcontext.WorkerErrorContext{
		Guild:   e.GuildId,
		User:    e.Author.Id,
		Channel: e.ChannelId,
	}
}
