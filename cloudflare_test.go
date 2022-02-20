package main

import (
	"os"
	"testing"
)

func TestCreateCfClientWithOutEnv(t *testing.T) {
	_, err := createCfClient()
	if err == nil {
		t.Fail()
	}
}

func TestCreateCfClientWithEnv(t *testing.T) {
	os.Setenv("CF_EMAIL", "email@email.com")
	os.Setenv("CF_TOKEN", "myAPIToken")
	os.Setenv("CF_ZONE", "myDNSZone")
	_, err := createCfClient()
	if err != nil {
		t.Fail()
	}
}
