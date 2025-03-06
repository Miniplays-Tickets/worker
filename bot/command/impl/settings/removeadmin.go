package settings

import (
	"fmt"
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/context"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/dbclient"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/utils"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	permcache "github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/channel/embed"
	"github.com/rxdn/gdl/objects/interaction"
	"github.com/rxdn/gdl/permission"
)

type RemoveAdminCommand struct{}

func (RemoveAdminCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "removeadmin",
		Description:     i18n.HelpRemoveAdmin,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permcache.Admin,
		Category:        command.Settings,
		Arguments: command.Arguments(
			command.NewRequiredArgument("user_or_role", "Bentuzer oder Rolle der Administrator Rechte für diesen Bot auf diesem Server entfernt werden sollen", interaction.OptionTypeMentionable, i18n.MessageRemoveAdminNoMembers),
		),
		DefaultEphemeral: true,
		Timeout:          time.Second * 5,
	}
}

func (c RemoveAdminCommand) GetExecutor() interface{} {
	return c.Execute
}

// TODO: Remove from existing tickets
func (c RemoveAdminCommand) Execute(ctx registry.CommandContext, id uint64) {
	usageEmbed := embed.EmbedField{
		Name:   "Usage",
		Value:  "`/removeadmin @User`\n`/removeadmin @Role`",
		Inline: false,
	}

	settings, err := ctx.Settings()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	mentionableType, valid := context.DetermineMentionableType(ctx, id)
	if !valid {
		ctx.ReplyWithFields(customisation.Red, i18n.Error, i18n.MessageRemoveAdminNoMembers, utils.ToSlice(usageEmbed))
		return
	}

	if mentionableType == context.MentionableTypeUser {
		// get guild object
		guild, err := ctx.Worker().GetGuild(ctx.GuildId())
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if guild.OwnerId == id {
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageOwnerMustBeAdmin)
			return
		}

		if ctx.UserId() == id {
			ctx.Reply(customisation.Red, i18n.Error, i18n.MessageRemoveStaffSelf)
			return
		}

		if err := dbclient.Client.Permissions.RemoveAdmin(ctx, ctx.GuildId(), id); err != nil {
			ctx.HandleError(err)
			return
		}

		if err := utils.ToRetriever(ctx.Worker()).Cache().SetCachedPermissionLevel(ctx, ctx.GuildId(), id, permcache.Support); err != nil {
			ctx.HandleError(err)
			return
		}
	} else if mentionableType == context.MentionableTypeRole {
		if err := dbclient.Client.RolePermissions.RemoveAdmin(ctx, ctx.GuildId(), id); err != nil {
			ctx.HandleError(err)
			return
		}

		if err := utils.ToRetriever(ctx.Worker()).Cache().SetCachedPermissionLevel(ctx, ctx.GuildId(), id, permcache.Support); err != nil {
			ctx.HandleError(err)
			return
		}
	} else {
		ctx.HandleError(fmt.Errorf("infallible"))
		return
	}

	ctx.Reply(customisation.Green, i18n.TitleRemoveAdmin, i18n.MessageRemoveAdminSuccess)

	// Remove user / role from thread notification channel
	if settings.TicketNotificationChannel != nil {
		_ = ctx.Worker().EditChannelPermissions(*settings.TicketNotificationChannel, channel.PermissionOverwrite{
			Id:    id,
			Type:  mentionableType.OverwriteType(),
			Allow: 0,
			Deny:  permission.BuildPermissions(permission.ViewChannel),
		})
	}
}
