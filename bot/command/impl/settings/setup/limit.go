package setup

import (
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/dbclient"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type LimitSetupCommand struct{}

func (LimitSetupCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "limit",
		Description:     i18n.HelpSetup,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Admin,
		Category:        command.Settings,
		Arguments: command.Arguments(
			command.NewRequiredArgument("limit", "Die Anzahl der maximalen Tickets die ein User gleichzeitig offen haben kann", interaction.OptionTypeInteger, i18n.SetupLimitInvalid),
		),
		Timeout: time.Second * 3,
	}
}

func (c LimitSetupCommand) GetExecutor() interface{} {
	return c.Execute
}

func (LimitSetupCommand) Execute(ctx registry.CommandContext, limit int) {
	if limit < 1 || limit > 10 {
		ctx.Reply(customisation.Red, i18n.TitleSetup, i18n.SetupLimitInvalid)
		return
	}

	if err := dbclient.Client.TicketLimit.Set(ctx, ctx.GuildId(), uint8(limit)); err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Reply(customisation.Green, i18n.TitleSetup, i18n.SetupLimitComplete, limit)
}
