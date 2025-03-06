package admin

import (
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type AdminCommand struct {
}

func (AdminCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "admin",
		Description:     i18n.HelpAdmin,
		Type:            interaction.ApplicationCommandTypeChatInput,
		Aliases:         []string{"a"},
		PermissionLevel: permission.Everyone,
		Children: []registry.Command{
			AdminBlacklistCommand{},
			AdminCheckBlacklistCommand{},
			AdminCheckPremiumCommand{},
			AdminGenPremiumCommand{},
			AdminGetOwnerCommand{},
			AdminListGuildEntitlementsCommand{},
			AdminListUserEntitlementsCommand{},
			AdminRecacheCommand{},
			AdminWhitelabelAssignGuildCommand{},
			AdminWhitelabelDataCommand{},
		},
		Category:   command.Settings,
		HelperOnly: true,
	}
}

func (c AdminCommand) GetExecutor() interface{} {
	return c.Execute
}

func (AdminCommand) Execute(_ registry.CommandContext) {
	// Cannot execute parent command directly
}
