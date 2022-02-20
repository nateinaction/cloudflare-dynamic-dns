package config

import (
	"encoding/json"
)

type Config struct {
	Records []Record `json:"records"`
}

type Record struct {
	Domain string `json:"domain"`
	ZoneId string `json:"zone_id"`
	Proxy  bool   `json:"proxy,omitempty"`
	Ttl    int    `json:"ttl,omitempty"`
}

func NewConfig(cfg []byte) (*Config, error) {
	config := &Config{}
	err := json.Unmarshal(cfg, &config)
	return config, err
}

// GetZones returns a list of zones and associated records
func (c *Config) GetZones() []string {
	zones := []string{}
	seen := map[string]bool{}
	for _, r := range c.Records {
		if _, ok := seen[r.ZoneId]; !ok {
			zones = append(zones, r.ZoneId)
			seen[r.ZoneId] = true
		}
	}
	return zones
}

// GetRecords returns a list of records
func (c *Config) GetRecords() map[string]Record {
	records := map[string]Record{}
	for _, r := range c.Records {
		records[r.Domain] = r
	}
	return records
}
