package network_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/network"
)

func TestIsIp(t *testing.T) {
	ip := &network.Ip{
		V4: "123",
	}
	if ip.IsIp() == nil {
		t.Error("Expected error, got nil")
	}

	ip2 := &network.Ip{
		V4: "10.0.0.1",
	}
	if ip2.IsIp() != nil {
		t.Error("Expected nil, got error")
	}
}

func TestMatch(t *testing.T) {
	ip1 := &network.Ip{
		V4: "123",
	}
	ip2 := &network.Ip{
		V4: "1234",
	}
	if ip1.Match(ip2) {
		t.Errorf("expected ip1 & ip2 mismatch, got match")
	}

	ip3 := &network.Ip{
		V4: "123",
	}
	if !ip1.Match(ip3) {
		t.Errorf("expected ip1 & ip3 match, got mismatch")
	}
}
