package settings

import (
	"time"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/bot/customisation"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type PanelCommand struct {
}

func (PanelCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:             "panel",
		Description:      i18n.HelpPanel,
		Type:             interaction.ApplicationCommandTypeChatInput,
		PermissionLevel:  permission.Admin,
		Category:         command.Settings,
		DefaultEphemeral: true,
		Timeout:          time.Second * 3,
	}
}

func (c PanelCommand) GetExecutor() interface{} {
	return c.Execute
}

func (PanelCommand) Execute(ctx registry.CommandContext) {
	ctx.Reply(customisation.Green, i18n.TitlePanel, i18n.MessagePanel, ctx.GuildId())
}
