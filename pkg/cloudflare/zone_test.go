package cloudflare_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
)

func TestNewZone(t *testing.T) {
	z, err := cloudflare.NewZone("12345")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if z.Id != "12345" {
		t.Errorf("expected Id to be 12345, got %v", z.Id)
	}

	_, err = cloudflare.NewZone("")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUrl(t *testing.T) {
	z, _ := cloudflare.NewZone("12345")
	expected := "https://api.cloudflare.com/client/v4/zones/12345/dns_records"
	if z.Url() != expected {
		t.Errorf("expected %s, got %v", expected, z.Url())
	}
}
