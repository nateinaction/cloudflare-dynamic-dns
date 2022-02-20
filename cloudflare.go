package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type cfClient struct {
	Email  string
	Token  string
	Zone   string
	Url    string
	Client *http.Client
}

// createCfClient Reads credentials from environment variables
func createCfClient() (*cfClient, error) {
	email := os.Getenv("CF_EMAIL")
	token := os.Getenv("CF_TOKEN")
	zone := os.Getenv("CF_ZONE")
	if email == "" || token == "" || zone == "" {
		return nil, fmt.Errorf("CF_EMAIL, CF_TOKEN, and CF_ZONE must be set")
	}
	return &cfClient{
		Email:  email,
		Token:  token,
		Url:    fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zone),
		Client: &http.Client{},
	}, nil
}

type DnsRecord struct {
	Id      string `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Name    string `json:"name,omitempty"`
	Proxied bool   `json:"proxied,omitempty"`
	Content string `json:"content,omitempty"`
	Ttl     int    `json:"ttl,omitempty"`
}

type DnsRecordsResponse struct {
	Records []DnsRecord `json:"result"`
}

func (c *cfClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
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
func (c *cfClient) GetDnsRecords() (*DnsRecordsResponse, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("%s?type=A", c.Url), nil)
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

	data := &DnsRecordsResponse{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}
	return data, nil
}

// UpdateRecord Sets the IP for an existing DNS record
func (c *cfClient) UpdateRecord(ip, id string) error {
	record := DnsRecord{
		Content: ip,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PATCH", fmt.Sprintf("%s/%s", c.Url, id), bytes.NewBuffer(data))
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
func (c *cfClient) CreateRecord(name, ip string) error {
	record := DnsRecord{
		Type:    "A",
		Name:    name,
		Content: ip,
		Ttl:     1,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", c.Url, bytes.NewBuffer(data))
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

// getPublicIP Finds the current network's public IP address
func getPublicIP() (string, error) {
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

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	recordsToUpdate := arrayFlags{}
	flag.Var(&recordsToUpdate, "r", "Records (sub.domain.tld) to update with current IP")
	flag.Parse()

	cfClient, err := createCfClient()
	if err != nil {
		log.Fatalln(err)
	}

	previousIp := "startup"
	for {
		publicIp, err := getPublicIP()
		if err != nil {
			log.Printf("failed to get public ip: %s\n", err)
		}

		if previousIp != publicIp {
			log.Printf("ip change detected: %s -> %s\n", previousIp, publicIp)
			cfRecords, err := cfClient.GetDnsRecords()
			if err != nil {
				log.Printf("failed to get cloudflare dns records: %s\n", err)
			}

			rMap := map[string]DnsRecord{}
			for _, record := range cfRecords.Records {
				rMap[record.Name] = record
			}

			for _, name := range recordsToUpdate {
				if _, ok := rMap[name]; !ok {
					if err := cfClient.CreateRecord(name, publicIp); err != nil {
						log.Printf("failed to create record: %s\n", err)
					}
					log.Printf("added record: %s\n", name)
				} else {
					// Prevent additional calls to cloudflare if the IP hasn't changed, for example, on application startup
					if rMap[name].Content != publicIp {
						if err := cfClient.UpdateRecord(publicIp, rMap[name].Id); err != nil {
							log.Printf("failed to update record: %s\n", err)
						}
						log.Printf("updated record: %s\n", name)
					}
				}
			}
			previousIp = publicIp
		}
		time.Sleep(3600 * time.Second)
	}
}
