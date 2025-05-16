package admin

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/bot/command/registry"
	"github.com/Miniplays-Tickets/worker/bot/customisation"
	"github.com/Miniplays-Tickets/worker/bot/dbclient"
	"github.com/Miniplays-Tickets/worker/bot/utils"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/rxdn/gdl/objects/interaction"
)

type AdminCheckBlacklistCommand struct {
}

func (AdminCheckBlacklistCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:            "check-blacklist",
		Description:     i18n.HelpAdmin,
		Type:            interaction.ApplicationCommandTypeChatInput,
		PermissionLevel: permission.Everyone,
		Category:        command.Settings,
		HelperOnly:      true,
		Arguments: command.Arguments(
			command.NewRequiredArgument("guild_id", "ID of the guild to unblacklist", interaction.OptionTypeString, i18n.MessageInvalidArgument),
		),
		Timeout: time.Second * 10,
	}
}

func (c AdminCheckBlacklistCommand) GetExecutor() interface{} {
	return c.Execute
}

func (AdminCheckBlacklistCommand) Execute(ctx registry.CommandContext, raw string) {
	guildId, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		ctx.ReplyRaw(customisation.Red, ctx.GetMessage(i18n.Error), "Invalid guild ID provided")
		return
	}

	isBlacklisted, reason, err := dbclient.Client.ServerBlacklist.IsBlacklisted(ctx, guildId)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if isBlacklisted {
		reasonFormatted := utils.ValueOrDefault(reason, "No reason provided")
		ctx.ReplyRaw(customisation.Orange, "Blacklist Check", fmt.Sprintf("This guild is blacklisted.\n```%s```", reasonFormatted))
	} else {
		ctx.ReplyRaw(customisation.Green, "Blacklist Check", "This guild is not blacklisted")
	}
}
