package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/config"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/secret"
)

type Client struct {
	Email       string
	Token       string
	UrlTemplate string
	Client      *http.Client
}

// createCfClient Reads credentials from environment variables
func NewClient(scrt *secret.Secret) (*Client, error) {
	return &Client{
		Email:       scrt.Email,
		Token:       scrt.Token,
		UrlTemplate: "https://api.cloudflare.com/client/v4/zones/%s/dns_records",
		Client:      &http.Client{},
	}, nil
}

type CfRecord struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
	Content string `json:"content,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
}

type CfRecordsResponse struct {
	Records []CfRecord `json:"result"`
}

func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Email", c.Email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-type", "application/json")
	return req, nil
}

func (c *Client) url(zone string) string {
	return fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zone)
}

// GetDnsRecords Retrieves DNS record data from Cloudflare
func (c *Client) GetDnsRecords(zones []string) ([]CfRecord, error) {
	results := []CfRecord{}
	for _, zone := range zones {
		resp, err := c.GetZoneRecords(zone)
		if err != nil {
			return nil, err
		}
		results = append(results, resp.Records...)
	}
	return results, nil
}

func (c *Client) GetZoneRecords(zone string) (*CfRecordsResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s?type=A", c.url(zone)), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &CfRecordsResponse{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateRecord Sets the IP for an existing DNS record
func (c *Client) UpdateRecord(r config.Record, id, ip string) error {
	ttl := 1
	if r.Ttl != 0 {
		ttl = r.Ttl
	}

	record := CfRecord{
		Type:    "A",
		Name:    r.Domain,
		Content: ip,
		Ttl:     ttl,
		Proxied: r.Proxy,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PUT", fmt.Sprintf("%s/%s", c.url(r.ZoneId), id), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// CreateRecord Sets the IP for a new DNS record
func (c *Client) CreateRecord(r config.Record, ip string) error {
	ttl := 1
	if r.Ttl != 0 {
		ttl = r.Ttl
	}

	record := CfRecord{
		Type:    "A",
		Name:    r.Domain,
		Content: ip,
		Ttl:     ttl,
		Proxied: r.Proxy,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", c.url(r.ZoneId), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
