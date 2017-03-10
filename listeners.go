package proxyprotocol

import (
	"net"

	proxyproto "github.com/armon/go-proxyproto"
	"github.com/mholt/caddy"
)

type Configs []Config

type Listener struct {
	caddy.Listener
	Configs []Config
}

func (c Configs) NewListener(l caddy.Listener) caddy.Listener {
	ln := &Listener{
		Listener: l,
		Configs:  []Config(c),
	}
	return ln
}

func (l *Listener) Accept() (net.Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	addr, ok := c.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return c, nil
	}
	for _, cfg := range l.Configs {
		for _, s := range cfg.Subnets {
			if s.Contains(addr.IP) {
				return proxyproto.NewConn(c, cfg.Timeout), nil
			}
		}
	}
	return c, nil
}
