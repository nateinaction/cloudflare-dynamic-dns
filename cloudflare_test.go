package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

func TestJsonify(t *testing.T) {
	var recordType, name, ip string = "A", "sub.example.com", "127.0.0.1"
	var proxied bool = true
	jsonData := jsonify(recordType, name, ip, proxied)
	if len(jsonData) == 0 {
		t.Fail()
	}
}

func TestGetCredentials(t *testing.T) {
	os.Setenv("CF_EMAIL", "email@email.com")
	os.Setenv("CF_KEY", "myAPIKey")
	os.Setenv("CF_ZONE", "myDNSZone")
	email, gapik, zone := getCredentials()

	if email == "" || gapik == "" || zone == "" {
		t.Fail()
	}
}
