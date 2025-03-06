package context

import (
	"context"
	"sync"
	"time"

	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/command/registry"
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/dbclient"
	"github.com/TicketsBot-cloud/database"
)

type StateCache struct {
	ctx registry.CommandContext

	settings   *database.Settings
	settingsMu sync.Mutex
}

func NewStateCache(ctx registry.CommandContext) *StateCache {
	return &StateCache{
		ctx: ctx,
	}
}

func (s *StateCache) Settings() (database.Settings, error) {
	s.settingsMu.Lock()
	defer s.settingsMu.Unlock()

	if s.settings != nil {
		return *s.settings, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	settings, err := dbclient.Client.Settings.Get(ctx, s.ctx.GuildId())
	if err != nil {
		return database.Settings{}, err
	}

	s.settings = &settings
	return settings, nil
}
