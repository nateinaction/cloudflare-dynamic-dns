package cloudflare_test

import (
	"fmt"
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/secret"
)

func TestNewRequest(t *testing.T) {
	scrt := &secret.Secret{
		Email: "blah",
		Token: "poke",
	}
	cf := cloudflare.NewClient(scrt)
	req, err := cf.NewRequest("POST", "someurl", nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	hEmail := req.Header.Get("X-Auth-Email")
	if hEmail != scrt.Email {
		t.Errorf("expected X-Auth-Email set to %s, got %v", scrt.Email, hEmail)
	}

	hAuth := req.Header.Get("Authorization")
	bToken := fmt.Sprintf("Bearer %s", scrt.Token)
	if hAuth != bToken {
		t.Errorf("expected Authorization set to %s, got %v", bToken, hAuth)
	}

	hContType := req.Header.Get("Content-Type")
	if hContType != "application/json" {
		t.Errorf("expected Content-Type set to application/json, got %v", hContType)
	}
}
