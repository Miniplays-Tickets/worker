package handlers

import (
	"errors"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/button/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/button/registry/matcher"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/context"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/constants"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/dbclient"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/logic"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	"github.com/TicketsBot/common/sentry"
)

type MultiPanelHandler struct{}

func (h *MultiPanelHandler) Matcher() matcher.Matcher {
	return &matcher.SimpleMatcher{
		CustomId: "multipanel",
	}
}

func (h *MultiPanelHandler) Properties() registry.Properties {
	return registry.Properties{
		Flags:   registry.SumFlags(registry.GuildAllowed),
		Timeout: constants.TimeoutOpenTicket,
	}
}

func (h *MultiPanelHandler) Execute(ctx *context.SelectMenuContext) {
	if len(ctx.InteractionData.Values) == 0 {
		return
	}

	panelCustomId := ctx.InteractionData.Values[0]

	panel, ok, err := dbclient.Client.Panel.GetByCustomId(ctx, ctx.GuildId(), panelCustomId)
	if err != nil {
		sentry.Error(err) // TODO: Proper context
		return
	}

	if ok {
		// TODO: Log this
		if panel.GuildId != ctx.GuildId() {
			return
		}

		// blacklist check
		blacklisted, err := ctx.IsBlacklisted(ctx)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if blacklisted {
			ctx.Reply(customisation.Red, i18n.TitleBlacklisted, i18n.MessageBlacklisted)
			return
		}

		if panel.FormId == nil {
			_, _ = logic.OpenTicket(ctx.Context, ctx, &panel, panel.Title, nil)
		} else {
			form, ok, err := dbclient.Client.Forms.Get(ctx, *panel.FormId)
			if err != nil {
				ctx.HandleError(err)
				return
			}

			if !ok {
				ctx.HandleError(errors.New("Form not found"))
				return
			}

			inputs, err := dbclient.Client.FormInput.GetInputs(ctx, form.Id)
			if err != nil {
				ctx.HandleError(err)
				return
			}

			if len(inputs) == 0 { // Don't open a blank form
				_, _ = logic.OpenTicket(ctx.Context, ctx, &panel, panel.Title, nil)
			} else {
				modal := buildForm(panel, form, inputs)
				ctx.Modal(modal)
			}
		}
	}
}
