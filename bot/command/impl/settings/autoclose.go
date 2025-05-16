package settings

import (
	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type AutoCloseCommand struct {
}

func (AutoCloseCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "autoclose",
		Description:     i18n.HelpAutoClose,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Support,
		Category:        command.Settings,
		Children: []registry.Command{
			AutoCloseConfigureCommand{},
			AutoCloseExcludeCommand{},
		},
	}
}

func (c AutoCloseCommand) GetExecutor() interface{} {
	return c.Execute
}

func (AutoCloseCommand) Execute(ctx registry.CommandContext) {
	// Can't call a parent command
}
