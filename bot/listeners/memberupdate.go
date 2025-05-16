package listeners

import (
	"context"
	"time"

	"github.com/Miniplays-Tickets/worker"
	"github.com/Miniplays-Tickets/worker/bot/utils"
	"github.com/TicketsBot-cloud/common/sentry"
	"github.com/rxdn/gdl/gateway/payloads/events"
)

// Remove user permissions when they leave
func OnMemberUpdate(worker *worker.Context, e events.GuildMemberUpdate) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3) // TODO: Propagate context
	defer cancel()

	span := sentry.StartSpan(ctx, "OnMemberUpdate")
	defer span.Finish()

	if err := utils.ToRetriever(worker).Cache().DeleteCachedPermissionLevel(ctx, e.GuildId, e.User.Id); err != nil {
		sentry.Error(err)
	}
}
