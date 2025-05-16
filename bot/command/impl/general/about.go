package general

import (
	"time"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/bot/customisation"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type AboutCommand struct {
}

func (AboutCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:             "about",
		Description:      i18n.HelpAbout,
		Type:             interaction.ApplicationCommandTypeChatInput,
		PermissionLevel:  permission.Everyone,
		Category:         command.General,
		MainBotOnly:      true,
		DefaultEphemeral: true,
		Timeout:          time.Second * 3,
	}
}

func (c AboutCommand) GetExecutor() interface{} {
	return c.Execute
}

func (AboutCommand) Execute(ctx registry.CommandContext) {
	ctx.Reply(customisation.Green, i18n.TitleAbout, i18n.MessageAbout)
}
