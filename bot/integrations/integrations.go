package integrations

import (
	"github.com/Dev-Miniplays/Ticketsv2-worker/bot/redis"
	"github.com/Dev-Miniplays/Ticketsv2-worker/config"
	"github.com/TicketsBot/common/integrations/bloxlink"
	"github.com/TicketsBot/common/webproxy"
)

var (
	WebProxy    *webproxy.WebProxy
	SecureProxy *SecureProxyClient
	Bloxlink    *bloxlink.BloxlinkIntegration
)

func InitIntegrations() {
	WebProxy = webproxy.NewWebProxy(config.Conf.WebProxy.Url, config.Conf.WebProxy.AuthHeaderName, config.Conf.WebProxy.AuthHeaderValue)
	Bloxlink = bloxlink.NewBloxlinkIntegration(redis.Client, WebProxy, config.Conf.Integrations.BloxlinkApiKey)
	SecureProxy = NewSecureProxy(config.Conf.Integrations.SecureProxyUrl)
}
