package config_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/config"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/network"
)

const exampleConfigJson = `{
	"records": [
		{
			"name": "example.com",
			"zone_id": "example-zone-id"
		},
		{
			"name": "sub.example.net",
			"zone_id": "example-zone-id",
			"proxied": true,
			"ttl": 300
		}
	]
}`

func TestNewConfig(t *testing.T) {
	c, err := config.NewConfig([]byte(exampleConfigJson))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(c.Records) != 2 {
		t.Errorf("expected 2 records, got %d", len(c.Records))
	}

	if c.Records[0].Name != "example.com" {
		t.Errorf("expected example.com, got %s", c.Records[0].Name)
	}

	if c.Records[0].Proxied != false {
		t.Errorf("expected proxy false, got %v", c.Records[0].Proxied)
	}

	if c.Records[1].Proxied != true {
		t.Errorf("expected proxy true, got %v", c.Records[1].Proxied)
	}
}

func TestGetZones(t *testing.T) {
	c := config.Config{
		Records: []cloudflare.Record{
			{
				Name:   "example.com",
				ZoneId: "example-zone-id",
			},
			{
				Name:   "sub.example.net",
				ZoneId: "example-zone-id",
			},
			{
				Name:   "sub2.example.net",
				ZoneId: "example-zone-id-2",
			},
		},
	}

	zones := c.GetZones()
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}

	if _, ok := zones["example-zone-id"]; !ok {
		t.Errorf("expected zone with id example-zone-id, got %v", zones)
	}
}

func TestGetRecords(t *testing.T) {
	c := config.Config{
		Records: []cloudflare.Record{
			{
				Type:   "A",
				Name:   "example.com",
				ZoneId: "example-zone-id",
			},
			{
				Type:   "A",
				Name:   "sub.example.net",
				ZoneId: "example-zone-id",
				Ttl:    300,
			},
			{
				Type:   "AAAA",
				Name:   "sub2.example.net",
				ZoneId: "example-zone-id-2",
			},
		},
	}

	r := c.GetRecords(&network.Ip{V4: "fake-ipv4", V6: "fake-ipv6"})
	if len(r) != 3 {
		t.Errorf("expected 3 records, got %d", len(r))
	}

	if _, ok := r["example.com"]; !ok {
		t.Errorf("expected record for example.com, got none")
	}

	if r["example.com"].Ttl != 1 {
		t.Errorf("expected ttl 1, got %d", r["example.com"].Ttl)
	}

	if r["sub.example.net"].Ttl != 300 {
		t.Errorf("expected ttl 300, got %d", r["sub.example.net"].Ttl)
	}

	if r["example.com"].Content != "fake-ipv4" {
		t.Errorf("expected content fake-ip, got %s", r["example.com"].Content)
	}

	if r["sub2.example.net"].Content != "fake-ipv6" {
		t.Errorf("expected content fake-ip, got %s", r["sub2.example.net"].Content)
	}
}
