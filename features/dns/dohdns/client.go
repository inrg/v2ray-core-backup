package dohdns

import (
	"sync"
	"time"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns"

	doh "github.com/babolivier/go-doh-client"
)

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
	rch := make(chan net.IP)
	resolvedIPs := make([]net.IP, 0)

	// will wait for LookupA and LookupAAAA
	wg.Add(2)

	// resolver will re-use keep-alive connection to DOH server
	go func() {
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
		close(rch)
	}()

	timeout := time.After(c.timeout)
	for {
		select {
		case ip, ok := <-rch:
			if ok {
				resolvedIPs = append(resolvedIPs, ip)
			} else {
				goto EndForIP
			}
		case <-timeout:
			newError("DOH timed out: ", host).AtDebug().WriteToLog()
			goto EndForIP
		}
	}
EndForIP:

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
		r, _, err := c.resolver.LookupA(host)
		if err == nil {
			for _, ip := range r {
				ip := net.ParseIP(ip.IP4)
				if ip != nil {
					rch <- ip
				}
			}
		}
		close(rch)
	}()

	resolvedIPs := make([]net.IP, 0)
	timeout := time.After(c.timeout)
	for {
		select {
		case ip, ok := <-rch:
			if ok {
				resolvedIPs = append(resolvedIPs, ip)
			} else {
				goto EndFor4
			}
		case <-timeout:
			newError("DOH timed out: ", host).AtDebug().WriteToLog()
			goto EndFor4
		}
	}
EndFor4:

	if len(resolvedIPs) == 0 {
		return nil, dns.ErrEmptyResponse
	}

	return resolvedIPs, nil
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host string) ([]net.IP, error) {
	newError("LookupIPv6 using DOH: ", host).AtDebug().WriteToLog()

	rch := make(chan net.IP)

	go func() {
		r, _, err := c.resolver.LookupAAAA(host)
		if err == nil {
			for _, ip := range r {
				ip := net.ParseIP(ip.IP6)
				if ip != nil {
					rch <- ip
				}
			}
		}
		close(rch)
	}()

	resolvedIPs := make([]net.IP, 0)
	timeout := time.After(c.timeout)
	for {
		select {
		case ip, ok := <-rch:
			if ok {
				resolvedIPs = append(resolvedIPs, ip)
			} else {
				goto EndFor6
			}
		case <-timeout:
			newError("DOH timed out: ", host).AtDebug().WriteToLog()
			goto EndFor6
		}
	}
EndFor6:

	if len(resolvedIPs) == 0 {
		return nil, dns.ErrEmptyResponse
	}

	return resolvedIPs, nil
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
