package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/network"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/secret"
)

type Client struct {
	Email       string
	Token       string
	UrlTemplate string
	Client      *http.Client
}

// NewClient Instantiates a new Cloudflare client
func NewClient(scrt *secret.Secret) *Client {
	return &Client{
		Email:  scrt.Email,
		Token:  scrt.Token,
		Client: &http.Client{},
	}
}

func (c *Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Auth-Email", c.Email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-type", "application/json")

	return req, nil
}

// GetDnsRecords Retrieves DNS record data from Cloudflare
func (c *Client) GetRecords(zs map[string]Zone) (map[string]Record, error) {
	results := map[string]Record{}
	for _, z := range zs {
		records, err := c.GetRecord(z)
		if err != nil {
			return nil, err
		}

		for _, r := range records {
			results[r.Name] = r
		}
	}
	return results, nil
}

func (c *Client) GetRecord(z Zone) ([]Record, error) {
	req, err := c.NewRequest("GET", fmt.Sprintf("%s?type=A", z.Url()), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &RecordResp{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	return data.Records, nil
}

// UpdateRecord Sets the IP for an existing DNS record
func (c *Client) UpdateRecord(r Record, z Zone, ip *network.Ip) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := c.NewRequest("PATCH", fmt.Sprintf("%s/%s", z.Url(), r.Id), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// CreateRecord Sets the IP for a new DNS record
func (c *Client) CreateRecord(r Record, z Zone, ip *network.Ip) error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	req, err := c.NewRequest("POST", z.Url(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
