package dohdns

import (
	"os"
	"testing"
	"time"

	doh "github.com/babolivier/go-doh-client"
)

func init() {
	os.Setenv("GODEBUG", "tls13=1")
}

func TestClient_LookupIPv4(t *testing.T) {
	type fields struct {
		resolver doh.Resolver
	}
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"dohipv4_CF", fields{doh.Resolver{"1.1.1.1", doh.IN}}, args{"www.nosuch.whatooo"}, true},
		{"dohipv4_CF", fields{doh.Resolver{"1.1.1.1", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohipv4_CF", fields{doh.Resolver{"1.1.1.1", doh.IN}}, args{"www.facebook.com"}, false},
		{"dohipv4_CF", fields{doh.Resolver{"1.1.1.1", doh.IN}}, args{"www.twitter.com"}, false},
		{"dohipv4_Google", fields{doh.Resolver{"dns.google", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohipv4_Cloudflare", fields{doh.Resolver{"cloudflare-dns.com", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohipv4_Cloudflare", fields{doh.Resolver{"9.9.9.9", doh.IN}}, args{"www.v2ray.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				resolver: tt.fields.resolver,
				timeout:  3 * time.Second,
			}
			got, err := c.LookupIPv4(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LookupIPv4() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%#v", got)
			if len(got) == 0 && !tt.wantErr {
				t.Errorf("Client.LookupIPv4(%s) didn't resolve any IP", tt.args.host)
			}
		})
	}
}

func TestClient_LookupIPv6(t *testing.T) {
	type fields struct {
		resolver doh.Resolver
	}
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"dohipv6_IBM", fields{doh.Resolver{"9.9.9.9", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohipv6_Google", fields{doh.Resolver{"dns.google", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohipv6_Cloudflare", fields{doh.Resolver{"cloudflare-dns.com", doh.IN}}, args{"www.v2ray.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				resolver: tt.fields.resolver,
				timeout:  3 * time.Second,
			}
			got, err := c.LookupIPv6(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LookupIPv6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%#v", got)
			if len(got) == 0 && !tt.wantErr {
				t.Errorf("Client.LookupIPv6(%s) didn't resolve any IP", tt.args.host)
			}
		})
	}
}

func TestClient_LookupIP(t *testing.T) {
	type fields struct {
		resolver doh.Resolver
	}
	type args struct {
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"dohip_IBM", fields{doh.Resolver{"9.9.9.9", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohip_Google", fields{doh.Resolver{"dns.google", doh.IN}}, args{"www.v2ray.com"}, false},
		{"dohip_Cloudflare", fields{doh.Resolver{"cloudflare-dns.com", doh.IN}}, args{"www.v2ray.com"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				resolver: tt.fields.resolver,
				timeout:  3 * time.Second,
			}
			got, err := c.LookupIP(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.LookupIPv6() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%#v", got)
			// www.v2ray.com should have 2 ipv4 and 2 ipv6 address (cloudflare frontend)
			if len(got) == 0 && !tt.wantErr {
				t.Errorf("Client.LookupIPv6(%s) didn't resolve any IP", tt.args.host)
			}
		})
	}
}
