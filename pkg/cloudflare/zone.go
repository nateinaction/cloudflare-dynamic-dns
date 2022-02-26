package cloudflare

import (
	"fmt"
)

type Zone struct {
	Id string
}

func NewZone(id string) (*Zone, error) {
	if id == "" {
		return nil, fmt.Errorf("zone id is required")
	}
	return &Zone{
		Id: id,
	}, nil
}

func (z *Zone) Url() string {
	return fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", z.Id)
}
