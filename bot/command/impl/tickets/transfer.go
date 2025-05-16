package tickets

import (
	"fmt"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/bot/constants"
	"github.com/Miniplays-Tickets/worker/bot/customisation"
	"github.com/Miniplays-Tickets/worker/bot/dbclient"
	"github.com/Miniplays-Tickets/worker/bot/logic"
	"github.com/Miniplays-Tickets/worker/bot/utils"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/channel"
	"github.com/rxdn/gdl/objects/interaction"
)

type TransferCommand struct {
}

func (TransferCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "transfer",
		Description:     i18n.HelpTransfer,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Support,
		Category:        command.Tickets,
		Arguments: command.Arguments(
			command.NewRequiredArgument("user", "Supporter zu dem das Ticket übertragen werden soll", interaction.OptionTypeUser, i18n.MessageInvalidUser),
		),
		Timeout: constants.TimeoutOpenTicket,
	}
}

func (c TransferCommand) GetExecutor() interface{} {
	return c.Execute
}

func (TransferCommand) Execute(ctx registry.CommandContext, userId uint64) {
	// Get ticket struct
	ticket, err := dbclient.Client.Tickets.GetByChannelAndGuild(ctx, ctx.ChannelId(), ctx.GuildId())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	// Verify this is a ticket channel
	if ticket.UserId == 0 {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNotATicketChannel)
		return
	}

	// Check if thread
	ch, err := ctx.Channel()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if ch.Type == channel.ChannelTypeGuildPrivateThread {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageClaimThread)
		return
	}

	member, err := ctx.Worker().GetGuildMember(ctx.GuildId(), userId)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	permissionLevel, err := permission.GetPermissionLevel(ctx, utils.ToRetriever(ctx.Worker()), member, ctx.GuildId())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if permissionLevel < permission.Support {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageInvalidUser)
		return
	}

	if err := logic.ClaimTicket(ctx, ctx, ticket, userId); err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ReplyPermanent(customisation.Green, i18n.TitleClaim, i18n.MessageClaimed, fmt.Sprintf("<@%d>", userId))
}
