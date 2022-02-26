package config

import (
	"encoding/json"
	"log"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/network"
)

type Config struct {
	Records []cloudflare.Record `json:"records"`
}

func NewConfig(cfg []byte) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal(cfg, &config)
	return config, err
}

// GetZones returns a list of zones and associated records
func (c *Config) GetZones() map[string]cloudflare.Zone {
	zones := map[string]cloudflare.Zone{}
	for _, r := range c.Records {
		if _, ok := zones[r.ZoneId]; !ok {
			z, err := cloudflare.NewZone(r.ZoneId)
			if err != nil {
				log.Printf("error loading zone data: %v", err)
			}
			zones[r.ZoneId] = *z
		}
	}
	return zones
}

// GetRecords returns a list of records
func (c *Config) GetRecords(ip *network.Ip) map[string]cloudflare.Record {
	records := map[string]cloudflare.Record{}
	for _, r := range c.Records {
		if r.Type == "" || r.Name == "" || r.ZoneId == "" {
			log.Printf("missing type, name or zone_id: %v", r)
			continue
		}
		if r.Ttl < 60 {
			r.Ttl = 1
		}
		if r.Type == "A" {
			r.Content = ip.V4
		} else if r.Type == "AAAA" {
			r.Content = ip.V6
		}
		records[r.Name] = r
	}
	return records
}
