package meshname

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/miekg/dns"
)

func init() { plugin.Register("meshname", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'meshname'
	if c.NextArg() {
		return plugin.Error("meshname", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		dnsClient := new(dns.Client)
		dnsClient.Timeout = 5000000000 // increased 5 seconds timeout

		return Meshname{dnsClient: dnsClient}
	})

	return nil
}
