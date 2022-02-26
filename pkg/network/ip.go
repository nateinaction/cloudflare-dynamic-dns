package network

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
)

type Ip struct {
	V4 string
	V6 string
}

// NewIp Finds the current network's public IP address
func NewIp() (*Ip, error) {
	ip := &Ip{}
	resp, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return ip, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ip, err
	}

	ip.V4 = string(bytes.TrimSpace(body))
	if err := ip.IsIp(); err != nil {
		return ip, err
	}
	return ip, nil
}

func (ip *Ip) IsIp() error {
	if net.ParseIP(ip.V4) == nil {
		return fmt.Errorf("not an ip: %s", ip.V4)
	}
	return nil
}

func (ip *Ip) Match(ip2 *Ip) bool {
	return reflect.DeepEqual(ip, ip2)
}
