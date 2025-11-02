package general

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

type JumpToTopCommand struct {
}

func (JumpToTopCommand) Properties() registry.Properties {
	return registry.Properties{
		Name:             "jumptotop",
		Description:      i18n.HelpJumpToTop,
		Type:             interaction.ApplicationCommandTypeChatInput,
		PermissionLevel:  permission.Everyone,
		Category:         command.General,
		DefaultEphemeral: true,
		Timeout:          time.Second * 5,
	}
}

func (c JumpToTopCommand) GetExecutor() interface{} {
	return c.Execute
}

func (JumpToTopCommand) Execute(ctx registry.CommandContext) {
	ticket, err := dbclient.Client.Tickets.GetByChannelAndGuild(ctx, ctx.ChannelId(), ctx.GuildId())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if ticket.Id == 0 {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageNotATicketChannel)
		return
	}

	if ticket.WelcomeMessageId == nil {
		ctx.Reply(customisation.Red, i18n.Error, i18n.MessageJumpToTopNoWelcomeMessage)
		return
	}

	messageLink := fmt.Sprintf("https://discord.com/channels/%d/%d/%d", ctx.GuildId(), ctx.ChannelId(), *ticket.WelcomeMessageId)

	embed := utils.BuildEmbed(ctx, customisation.Green, i18n.TitleJumpToTop, i18n.MessageJumpToTopContent, nil)
	res := command.NewEphemeralEmbedMessageResponse(embed)
	res.Components = []component.Component{
		component.BuildActionRow(component.BuildButton(component.Button{
			Label:    ctx.GetMessage(i18n.ClickHere),
			Style:    component.ButtonStyleLink,
			Emoji:    nil,
			Url:      utils.Ptr(messageLink),
			Disabled: false,
		})),
	}

	if _, err := ctx.ReplyWith(res); err != nil {
		ctx.HandleError(err)
		return
	}
}
