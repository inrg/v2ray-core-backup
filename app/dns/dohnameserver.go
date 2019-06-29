// +build !confonly

package dns

import (
	"context"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns/dohdns"
)

type dohNameServer struct {
	*dohdns.Client
}

func (s *dohNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {
	if option.IPv4Enable && option.IPv6Enable {
		return s.Client.LookupIP(domain)
	}

	if option.IPv4Enable {
		return s.Client.LookupIPv4(domain)
	}

	if option.IPv6Enable {
		return s.Client.LookupIPv6(domain)
	}

	return nil, newError("neither IPv4 nor IPv6 is enabled")
}

func (s *dohNameServer) Name() string {
	return "dohdns"
}

func NewDOHNameServer(host string, timeout time.Duration) *dohNameServer {
	return &dohNameServer{
		dohdns.New(host, timeout),
	}
}
