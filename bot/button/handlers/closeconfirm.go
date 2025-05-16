package handlers

import (
	"github.com/Miniplays-Tickets/worker/bot/button/registry"
	"github.com/Miniplays-Tickets/worker/bot/button/registry/matcher"
	"github.com/Miniplays-Tickets/worker/bot/command/context"
	"github.com/Miniplays-Tickets/worker/bot/constants"
	"github.com/Miniplays-Tickets/worker/bot/logic"
)

type CloseConfirmHandler struct{}

func (h *CloseConfirmHandler) Matcher() matcher.Matcher {
	return &matcher.SimpleMatcher{
		CustomId: "close_confirm",
	}
}

func (h *CloseConfirmHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed),
		Timeout: constants.TimeoutCloseTicket,
	}
}

func (h *CloseConfirmHandler) Execute(ctx *context.ButtonContext) {
	// TODO: IntoPanelContext()?
	logic.CloseTicket(ctx.Context, ctx, nil, false)
}
