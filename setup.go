package altonions

import (
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	proxyprotocol "github.com/mastercactapus/caddy-proxyprotocol"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/header"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type Config struct {
	Onions  []string
	MaxAge  uint32
	Persist bool
}

func init() {
	caddy.RegisterPlugin("altonions", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	var config Config
	var proxyConfigs []proxyprotocol.Config
	var err error

	hsv3, _ := regexp.Compile("[a-z2-8]{56}\\.onion:[0-9]{1,5}")

	for c.Next() {
		var proxyConfig proxyprotocol.Config
		for c.NextArg() {
			onion := c.Val()
			if !hsv3.MatchString(onion) {
				return c.ArgErr()
			}
			config.Onions = append(config.Onions, onion)
		}

		if c.NextBlock() {
			switch c.Val() {
			case "ma":
				if !c.NextArg() {
					return c.ArgErr()
				}
				seconds, err := strconv.ParseUint(c.Val(), 10, 32)
				if err != nil {
					return err
				}
				config.MaxAge = uint32(seconds)
			case "persist":
				if !c.NextArg() {
					return c.ArgErr()
				}
				config.Persist, err = strconv.ParseBool(c.Val())
				if err != nil {
					return err
				}

			case "subnets":
				if !c.NextArg() {
					return c.ArgErr()
				}
				_, n, err := net.ParseCIDR(c.Val())
				if err != nil {
					return err
				}
				proxyConfig.Subnets = append(proxyConfig.Subnets, n)
			case "timeout":
				if !c.NextArg() {
					return c.ArgErr()
				}
				proxyConfig.Timeout, err = time.ParseDuration(c.Val())
				if err != nil {
					return err
				}
			default:
				return c.ArgErr()
			}
		}

		if proxyConfig.Subnets == nil {
			_, cidr, _ := net.ParseCIDR("::1/128")
			proxyConfig.Subnets = append(proxyConfig.Subnets, cidr)
			_, cidr, _ = net.ParseCIDR("127.0.0.1/32")
			proxyConfig.Subnets = append(proxyConfig.Subnets, cidr)
		}
		proxyConfigs = append(proxyConfigs, proxyConfig)

		if c.NextBlock() {
			return c.ArgErr()
		}
		if c.Next() {
			return c.ArgErr()
		}
	}

	if proxyConfigs == nil {
		return c.ArgErr()
	}

	// Generate the Alt-Svc header
	var altOnions []string
	for _, onion := range config.Onions {
		value := "h2=\"" + onion + "\""
		if config.MaxAge != 0 {
			value = value + "; ma=" + strconv.FormatUint(uint64(config.MaxAge), 10)
		}
		if config.Persist {
			value = value + "; persist=1"
		}
		altOnions = append(altOnions, value)
	}

	// Setup the Alt-Svc header middleware
	var head header.Rule
	head.Headers = http.Header{}
	head.Path = "/"
	head.Headers.Add("Alt-Svc", strings.Join(altOnions[:], ","))

	httpserver.GetConfig(c).AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
		return header.Headers{Next: next, Rules: []header.Rule{head}}
	})

	// Setup the proxy protocol listener middleware
	proxylistener := proxyprotocol.Configs(proxyConfigs).NewListener
	httpserver.GetConfig(c).AddListenerMiddleware(proxylistener)

	return nil
}
