package admin

import (
	"fmt"
	"time"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/bot/customisation"
	"github.com/Miniplays-Tickets/worker/bot/dbclient"
	"github.com/Miniplays-Tickets/worker/bot/utils"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
	"github.com/TicketsBot-cloud/gdl/objects/interaction/component"
)

type AdminListUserEntitlementsCommand struct {
}

func (AdminListUserEntitlementsCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "list-user-entitlements",
		Description:     i18n.HelpAdmin,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Everyone,
		Category:        command.Settings,
		HelperOnly:      true,
		Arguments: command.Arguments(
			command.NewRequiredArgument("user", "User to fetch entitlements for", interaction.OptionTypeUser, i18n.MessageInvalidArgument),
		),
		Timeout: time.Second * 15,
	}
}

func (c AdminListUserEntitlementsCommand) GetExecutor() any {
	return c.Execute
}

func (AdminListUserEntitlementsCommand) Execute(ctx registry.CommandContext, userId uint64) {
	// List entitlements that have expired in the past 30 days
	entitlements, err := dbclient.Client.Entitlements.ListUserSubscriptions(ctx, userId, time.Hour*24*30)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if len(entitlements) == 0 {
		ctx.ReplyRaw(customisation.Red, ctx.GetMessage(i18n.Error), "This user has no entitlements")
		return
	}

	innerComponents := []component.Component{}

	for _, entitlement := range entitlements {
		sourceFormatted := string(entitlement.Source)
		if sourceFormatted == "" {
			sourceFormatted = "None"
		}
		value := fmt.Sprintf(
			"#### %s\n\n**Tier:** %s\n**Source:** %s\n**Expires:** <t:%d>\n**SKU ID:** %s\n**SKU Priority:** %d\n\n",
			entitlement.SkuLabel,
			entitlement.Tier,
			sourceFormatted,
			entitlement.ExpiresAt.Unix(),
			entitlement.SkuId.String(),
			entitlement.SkuPriority,
		)

		innerComponents = append(innerComponents, component.BuildTextDisplay(component.TextDisplay{Content: value}))
	}

	ctx.ReplyWith(command.NewMessageResponseWithComponents([]component.Component{
		utils.BuildContainerWithComponents(
			ctx,
			customisation.Orange,
			"Admin - User Entitlements",
			innerComponents,
		),
	}))
}
