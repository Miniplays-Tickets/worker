package handlers

import (
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/button"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/button/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/button/registry/matcher"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/context"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	prem "github.com/Dev-Miniplays/Ticketsv2-worker/bot/premium"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	"github.com/TicketsBot/common/permission"
)

type PremiumKeyButtonHandler struct{}

func (h *PremiumKeyButtonHandler) Matcher() matcher.Matcher {
	return &matcher.SimpleMatcher{
		CustomId: "open_premium_key_modal",
	}
}

func (h *PremiumKeyButtonHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed),
		Timeout: time.Second * 3,
	}
}

func (h *PremiumKeyButtonHandler) Execute(ctx *context.ButtonContext) {
	// Get permission level
	permissionLevel, err := ctx.UserPermissionLevel(ctx)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if permissionLevel < permission.Admin {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNoPermission)
		return
	}

	ctx.Modal(button.ResponseModal{
		Data: prem.BuildKeyModal(ctx.GuildId()),
	})
}
