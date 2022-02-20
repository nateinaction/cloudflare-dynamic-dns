package secret

import (
	"encoding/json"
)

type Secret struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewSecret(scrt []byte) (*Secret, error) {
	secret := &Secret{}
	err := json.Unmarshal(scrt, &secret)
	return secret, err
}
