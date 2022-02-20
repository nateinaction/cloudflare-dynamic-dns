package publicip

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// Lookup Finds the current network's public IP address
func Lookup() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes.TrimSpace(body)), nil
}
