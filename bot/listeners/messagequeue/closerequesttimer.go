package messagequeue

import (
	"context"

	"github.com/Miniplays-Tickets/worker/bot/cache"
	cmdcontext "github.com/Miniplays-Tickets/worker/bot/command/context"
	"github.com/Miniplays-Tickets/worker/bot/constants"
	"github.com/Miniplays-Tickets/worker/bot/dbclient"
	"github.com/Miniplays-Tickets/worker/bot/logic"
	"github.com/Miniplays-Tickets/worker/bot/metrics/statsd"
	"github.com/Miniplays-Tickets/worker/bot/redis"
	"github.com/Miniplays-Tickets/worker/bot/utils"
	"github.com/TicketsBot-cloud/common/closerequest"
	"github.com/TicketsBot-cloud/common/sentry"
	"github.com/TicketsBot-cloud/database"
	"go.uber.org/zap"
)

func ListenCloseRequestTimer(logger *zap.Logger) {
	ch := make(chan database.CloseRequest)
	go closerequest.Listen(redis.Client, ch)

	for request := range ch {
		statsd.Client.IncrementKey(statsd.AutoClose)

		request := request
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), constants.TimeoutCloseTicket)
			defer cancel()

			logger.Debug("Processing close request",
				zap.Int("ticket_id", request.TicketId),
				zap.Uint64("guild_id", request.GuildId),
				zap.Uint64("user_id", request.UserId),
			)

			// get ticket
			ticket, err := dbclient.Client.Tickets.Get(ctx, request.TicketId, request.GuildId)
			if err != nil {
				logger.Error("Failed to fetch ticket",
					zap.Int("ticket_id", request.TicketId),
					zap.Uint64("guild_id", request.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			// get worker
			worker, err := buildContext(ctx, ticket, cache.Client)
			if err != nil {
				logger.Error("Failed to build worker context",
					zap.Int("ticket_id", request.TicketId),
					zap.Uint64("guild_id", request.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			// query already checks, but just to be sure
			if ticket.ChannelId == nil {
				logger.Warn("Ticket channel ID is nil",
					zap.Int("ticket_id", request.TicketId),
					zap.Uint64("guild_id", request.GuildId),
				)
				return
			}

			// get premium status
			premiumTier, err := utils.PremiumClient.GetTierByGuildId(ctx, ticket.GuildId, true, worker.Token, worker.RateLimiter)
			if err != nil {
				logger.Error("Failed to get premium tier",
					zap.Int("ticket_id", request.TicketId),
					zap.Uint64("guild_id", request.GuildId),
					zap.Error(err),
				)
				sentry.Error(err)
				return
			}

			cc := cmdcontext.NewAutoCloseContext(ctx, worker, ticket.GuildId, *ticket.ChannelId, request.UserId, premiumTier)
			logic.CloseTicket(ctx, cc, request.Reason, true)

			logger.Info("Successfully processed close request",
				zap.Int("ticket_id", request.TicketId),
				zap.Uint64("guild_id", request.GuildId),
				zap.Uint64("user_id", request.UserId),
			)
		}()
	}
}
