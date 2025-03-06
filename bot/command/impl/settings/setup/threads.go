package setup

import (
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/customisation"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/dbclient"
	"github.com/Dev-Miniplays/Ticketsv2-worker/i18n"
	"github.com/TicketsBot/common/permission"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/interaction"
)

type ThreadsSetupCommand struct{}

func (ThreadsSetupCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "use-threads",
		Description:     i18n.HelpSetup,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Admin,
		Category:        command.Settings,
		Arguments: command.Arguments(
			command.NewRequiredArgument("use_threads", "Ob Private Threads oder nicht für Tickets genutzt werden sollen", interaction.OptionTypeBoolean, "infallible"),
			command.NewOptionalArgument("ticket_notification_channel", "Der Kanal wo Benarichtigungen für offene Tickets gesendet werden sollen", interaction.OptionTypeChannel, "infallible"),
		),
		InteractionOnly: true,
		Timeout:         time.Second * 5,
	}
}

func (c ThreadsSetupCommand) GetExecutor() interface{} {
	return c.Execute
}

func (ThreadsSetupCommand) Execute(ctx registry.CommandContext, useThreads bool, channelId *uint64) {
	if useThreads {
		if channelId == nil {
			ctx.Reply(customisation.Red, i18n.Error, i18n.SetupThreadsNoNotificationChannel)
			return
		}

		ch, err := ctx.Worker().GetChannel(*channelId)
		if err != nil {
			ctx.HandleError(err)
			return
		}

		if ch.Type != channel.ChannelTypeGuildText {
			ctx.Reply(customisation.Red, i18n.Error, i18n.SetupThreadsNotificationChannelType)
			return
		}

		if err := dbclient.Client.Settings.EnableThreads(ctx, ctx.GuildId(), *channelId); err != nil {
			ctx.HandleError(err)
			return
		}

		ctx.Reply(customisation.Green, i18n.TitleSetup, i18n.SetupThreadsSuccess)
	} else {
		if err := dbclient.Client.Settings.DisableThreads(ctx, ctx.GuildId()); err != nil {
			ctx.HandleError(err)
			return
		}

		ctx.Reply(customisation.Green, i18n.TitleSetup, i18n.SetupThreadsDisabled)
	}
}
