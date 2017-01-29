package proxyprotocol

import (
	"net"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Config struct {
	Timeout time.Duration
	Subnets []*net.IPNet
}

func init() {
	caddy.RegisterPlugin("proxyprotocol", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	var configs []Config
	var err error

	for c.Next() {
		var cfg Config
		for c.NextArg() {
			_, n, err := net.ParseCIDR(c.Val())
			if err != nil {
				return err
			}
			cfg.Subnets = append(cfg.Subnets, n)
		}

		if c.NextBlock() {
			switch c.Val() {
			case "timeout":
				if !c.NextArg() {
					return c.ArgErr()
				}
				cfg.Timeout, err = time.ParseDuration(c.Val())
				if err != nil {
					return err
				}
			default:
				return c.ArgErr()
			}
		}
		configs = append(configs, cfg)
		if c.NextBlock() {
			return c.ArgErr()
		}
	}
	if configs != nil {
		httpserver.GetConfig(c).AddListenerMiddleware(Configs(configs).NewListener)
	}

	return nil
}
