package dohdns

import (
	"sync"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns"

	doh "github.com/babolivier/go-doh-client"
)

func panicHandle() {
	if r := recover(); r != nil {
		newError("dohdns panic handled: ", r).AtDebug().WriteToLog()
	}
}

// Client is an implementation of dns.Client, which queries localhost for DNS.
type Client struct {
	resolver doh.Resolver
	timeout  time.Duration
}

// Type implements common.HasType.
func (*Client) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (*Client) Start() error { return nil }

// Close implements common.Closable.
func (*Client) Close() error { return nil }

// LookupIP implements Client.
func (c *Client) LookupIP(host string) ([]net.IP, error) {
	newError("LookupIP using DOH: ", host).AtDebug().WriteToLog()

	var wg sync.WaitGroup
	var co sync.Once
	rch := make(chan net.IP)
	resolvedIPs := make([]net.IP, 0)

	// will wait for LookupA and LookupAAAA
	wg.Add(2)

	// resolver will re-use keep-alive connection to DOH server
	go func() {
		defer panicHandle()
		defer wg.Done()
		rec, _, err := c.resolver.LookupA(host)
		if err == nil {
			for _, r := range rec {
				ip := net.ParseIP(r.IP4)
				if ip != nil {
					rch <- ip
				}
			}
		}
	}()

	go func() {
		defer panicHandle()
		defer wg.Done()
		rec, _, err := c.resolver.LookupAAAA(host)
		if err == nil {
			for _, r := range rec {
				ip := net.ParseIP(r.IP6)
				if ip != nil {
					rch <- ip
				}
			}
		}
	}()

	// wait for results
	go func() {
		wg.Wait()
		co.Do(func() {
			close(rch)
		})
	}()

	// timeout waiter, resolver using http.Client, has a timeout of fixed 30s
	go func() {
		<-time.After(c.timeout)
		co.Do(func() {
			close(rch)
		})
	}()

	for ip := range rch {
		resolvedIPs = append(resolvedIPs, ip)
	}

	if len(resolvedIPs) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return resolvedIPs, nil
}

// LookupIPv4 implements IPv4Lookup.
func (c *Client) LookupIPv4(host string) ([]net.IP, error) {
	newError("LookupIPv4 using DOH: ", host).AtDebug().WriteToLog()

	rch := make(chan net.IP)

	go func() {
		<-time.After(c.timeout)
		close(rch)
	}()

	go func() {
		defer panicHandle()
		r, _, err := c.resolver.LookupA(host)
		if err == nil {
			for _, ip := range r {
				ip := net.ParseIP(ip.IP4)
				if ip != nil {
					rch <- ip
				}
			}
		}
	}()

	ipv4 := make([]net.IP, 0)
	if ip, ok := <-rch; ok {
		ipv4 = append(ipv4, ip)
	} else {
		return nil, dns.ErrEmptyResponse
	}
	return ipv4, nil
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host string) ([]net.IP, error) {
	newError("LookupIPv6 using DOH: ", host).AtDebug().WriteToLog()

	rch := make(chan net.IP)

	go func() {
		<-time.After(c.timeout)
		close(rch)
	}()

	go func() {
		defer panicHandle()
		r, _, err := c.resolver.LookupAAAA(host)
		if err == nil {
			for _, ip := range r {
				ip := net.ParseIP(ip.IP6)
				if ip != nil {
					rch <- ip
				}
			}
		}
	}()

	ipv6 := make([]net.IP, 0)
	if ip, ok := <-rch; ok {
		ipv6 = append(ipv6, ip)
	} else {
		return nil, dns.ErrEmptyResponse
	}
	return ipv6, nil
}

// New create a new dns.Client
func New(host string, timeout time.Duration) *Client {
	return &Client{
		doh.Resolver{
			Host:  host,
			Class: doh.IN,
		},
		timeout,
	}
}
