package registry

import (
	"time"

	"github.com/Miniplays-Tickets/worker/bot/command"
	"github.com/Miniplays-Tickets/worker/i18n"
	"github.com/TicketsBot-cloud/common/permission"
	"github.com/TicketsBot-cloud/gdl/objects/interaction"
)

type Properties struct {
	Name             string
	Description      i18n.MessageId
	Type             interaction.ApplicationCommandType
	Aliases          []string
	PermissionLevel  permission.PermissionLevel
	Children         []Command // TODO: Map
	PremiumOnly      bool
	Category         command.Category
	AdminOnly        bool
	HelperOnly       bool
	InteractionOnly  bool
	MessageOnly      bool
	MainBotOnly      bool
	Arguments        []command.Argument
	DefaultEphemeral bool
	Timeout          time.Duration

	SetupFunc func()
}
