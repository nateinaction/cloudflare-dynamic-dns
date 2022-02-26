package cloudflare

import "reflect"

type Record struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
	Content string `json:"content,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
	ZoneId  string `json:"zone_id,omitempty"`
}

type RecordResp struct {
	Records []Record `json:"result"`
}

func (r Record) Match(r2 Record) bool {
	return reflect.DeepEqual(r, r2)
}
