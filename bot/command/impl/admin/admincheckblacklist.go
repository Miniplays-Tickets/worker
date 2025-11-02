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
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
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

func (c AdminCheckBlacklistCommand) GetExecutor() any {
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

	response := fmt.Sprintf("This guild is not blacklisted\n\n**Guild ID:** `%d`", guildId)

	if isBlacklisted {
		response = fmt.Sprintf("This guild is blacklisted.\n\n**Guild ID:** `%d`\n**Reason:** `%s`", guildId, utils.ValueOrDefault(reason, "No reason provided"))
	}

	ctx.ReplyWith(command.NewMessageResponseWithComponents(utils.Slice(
		utils.BuildContainerRaw(
			ctx,
			customisation.Orange,
			"Admin - Blacklist Check",
			response,
		),
	)))
}
