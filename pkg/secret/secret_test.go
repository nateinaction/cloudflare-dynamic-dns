package secret_test

import (
	"testing"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/secret"
)

/* #nosec G101 */
const exampleSecretJson = `{
	"email": "cloudflare_account_email",
	"token": "cloudflare_api_token"
}`

func TestNewSecret(t *testing.T) {
	s, err := secret.NewSecret([]byte(exampleSecretJson))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if s.Email != "cloudflare_account_email" {
		t.Errorf("expected cloudflare_account_email, got %s", s.Email)
	}

	if s.Token != "cloudflare_api_token" {
		t.Errorf("expected cloudflare_api_token, got %s", s.Token)
	}
}
