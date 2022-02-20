package config_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/config"
)

const exampleConfigJson = `{
	"records": [
		{
			"domain": "example.com",
			"zone_id": "example-zone-id"
		},
		{
			"domain": "sub.example.net",
			"zone_id": "example-zone-id",
			"proxy": true,
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

	if c.Records[0].Domain != "example.com" {
		t.Errorf("expected example.com, got %s", c.Records[0].Domain)
	}

	if c.Records[0].Proxy != false {
		t.Errorf("expected proxy false, got %v", c.Records[0].Proxy)
	}

	if c.Records[1].Proxy != true {
		t.Errorf("expected proxy true, got %v", c.Records[1].Proxy)
	}
}

func TestGetZones(t *testing.T) {
	c := config.Config{
		Records: []config.Record{
			{
				Domain: "example.com",
				ZoneId: "example-zone-id",
			},
			{
				Domain: "sub.example.net",
				ZoneId: "example-zone-id",
			},
			{
				Domain: "sub2.example.net",
				ZoneId: "example-zone-id-2",
			},
		},
	}

	zones := c.GetZones()
	if len(zones) != 2 {
		t.Errorf("expected 2 zones, got %d", len(zones))
	}
}

func TestGetRecords(t *testing.T) {
	c := config.Config{
		Records: []config.Record{
			{
				Domain: "example.com",
				ZoneId: "example-zone-id",
			},
			{
				Domain: "sub.example.net",
				ZoneId: "example-zone-id",
			},
			{
				Domain: "sub2.example.net",
				ZoneId: "example-zone-id-2",
			},
		},
	}

	r := c.GetRecords()
	if len(r) != 3 {
		t.Errorf("expected 3 records, got %d", len(r))
	}

	if _, ok := r["example.com"]; !ok {
		t.Errorf("expected record for example.com, got none")
	}
}
