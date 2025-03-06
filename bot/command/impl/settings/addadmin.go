package settings

import (
	"fmt"
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/context"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/utils"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	permcache "github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/objects/interaction/component"
)

type AddAdminCommand struct{}

func (AddAdminCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "addadmin",
		Description:     i18n.HelpAddAdmin,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permcache.Admin,
		Category:        command.Settings,
		InteractionOnly: true,
		Arguments: command.Arguments(
			command.NewRequiredArgument("user_or_role", "Bentuzer oder Rolle der Administrator Rechte für diesen Bot auf diesem Server gegeben werden sollen", interaction.OptionTypeMentionable, i18n.MessageAddAdminNoMembers),
		),
		DefaultEphemeral: true,
		Timeout:          time.Second * 3,
	}
}

func (c AddAdminCommand) GetExecutor() interface{} {
	return c.Execute
}

func (c AddAdminCommand) Execute(ctx registry.CommandContext, id uint64) {
	usageEmbed := embed.EmbedField{
		Name:   "Usage",
		Value:  "`/addadmin @User`\n`/addadmin @Role`",
		Inline: false,
	}

	mentionableType, valid := context.DetermineMentionableType(ctx, id)
	if !valid {
		ctx.ReplyWithFields(customisation.Red, i18n.Error, i18n.MessageAddSupportNoMembers, utils.ToSlice(usageEmbed))
		return
	}

	var mention string
	if mentionableType == context.MentionableTypeUser {
		mention = fmt.Sprintf("<@%d>", id)
	} else if mentionableType == context.MentionableTypeRole {
		mention = fmt.Sprintf("<@&%d>", id)
	} else {
		ctx.HandleError(fmt.Errorf("unknown mentionable type: %d", mentionableType))
		return
	}

	// Send confirmation message
	e := utils.BuildEmbed(ctx, customisation.Green, i18n.TitleAddAdmin, i18n.MessageAddAdminConfirm, nil, mention)
	res := command.NewEphemeralEmbedMessageResponseWithComponents(e, utils.Slice(component.BuildActionRow(
		component.BuildButton(component.Button{
			Label:    ctx.GetMessage(i18n.Confirm),
			CustomId: fmt.Sprintf("addadmin-%d-%d", mentionableType, id),
			Style:    component.ButtonStylePrimary,
			Emoji:    nil,
		}),
	)))

	if _, err := ctx.ReplyWith(res); err != nil {
		ctx.HandleError(err)
	}
}
