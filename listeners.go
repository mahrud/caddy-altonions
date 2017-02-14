package proxyprotocol

import (
	"net"
	"os"

	proxyproto "github.com/armon/go-proxyproto"
	"github.com/mholt/caddy"
)

type Configs []Config

type Listener struct {
	net.Listener
	Configs []Config
}
type CaddyListener struct {
	*Listener
}

func (c *CaddyListener) File() (*os.File, error) {
	return c.Listener.Listener.(caddy.Listener).File()
}

func (c Configs) NewListener(l net.Listener) net.Listener {
	ln := &Listener{
		Listener: l,
		Configs:  []Config(c),
	}
	if _, ok := l.(caddy.Listener); ok {
		return &CaddyListener{Listener: ln}
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
