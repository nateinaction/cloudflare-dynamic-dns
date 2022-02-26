package cloudflare_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
)

func TestMatch(t *testing.T) {
	r1 := cloudflare.Record{
		Id:      "abc",
		Type:    "def",
		Content: "123",
	}
	r2 := cloudflare.Record{
		Id:      "abc",
		Type:    "def",
		Content: "1234",
	}
	if r1.Match(r2) {
		t.Errorf("expected r1 & r2 mismatch , got match")
	}

	r3 := cloudflare.Record{
		Id:      "abc",
		Type:    "def",
		Content: "123",
	}
	if !r1.Match(r3) {
		t.Errorf("expected r1 & r3 match, got mismatch")
	}
}
