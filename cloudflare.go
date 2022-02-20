package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//Global variable declarations
var interval time.Duration = 30
var tableData []table = make([]table, 10)
var setTime []string = make([]string, 10)
var numOfRecords int
var email string
var token string
var zone string
var recordToUpdate arrayFlags

//For html table
type table struct {
	Name  string
	IP    string
	Proxy bool
	Time  string
}

//JSON response struct
type response struct {
	Result []struct {
		Identifier string `json:"id"`
		Type       string `json:"type"`
		Name       string `json:"name"`
		Proxied    bool   `json:"proxied"`
		Content    string `json:"content"`
	} `json:"result"`
}

//JSON PUT struct
type sendme struct {
	RecordType string `json:"type"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	Proxied    bool   `json:"proxied"`
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//Uniform http.NewRequest template for mutliple operations
func httpRequest(client *http.Client, reqType string, url string, instruction []byte, email string, token string, zone string) []byte {
	var req *http.Request
	var err error
	if instruction == nil {
		req, err = http.NewRequest(reqType, url, nil)
	} else {
		req, err = http.NewRequest(reqType, url, bytes.NewBuffer(instruction))
	}
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("X-Auth-Email", email)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return body
}

//Unmarshals the JSON response into the 'response' struct type
func unjsonify(body []byte) response {
	var jsonData response
	err := json.Unmarshal([]byte(body), &jsonData)
	if err != nil {
		log.Fatalln(err)
	}
	return jsonData
}

//Creates a JSON payload of type 'sendme'
func jsonify(recordType string, name string, ip string, proxied bool) []byte {
	data := sendme{
		RecordType: recordType,
		Name:       name,
		Content:    ip,
		Proxied:    proxied,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalln(err)
	}
	return jsonData
}

//Finds the computer's current public IP address and returns it
func getIP() string {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(bytes.TrimSpace(body))
}

//Reads credentials via shell variables
func getCredentials() (string, string, string) {
	email = os.Getenv("CF_EMAIL")
	token = os.Getenv("CF_KEY")
	zone = os.Getenv("CF_ZONE")

	//Checks to make sure email, global API key, and zone are set.
	if email == "" || token == "" || zone == "" {
		fmt.Println("Account email, API token, or zone is not set")
		os.Exit(1)
	}
	return email, token, zone
}

func update(client *http.Client, recordNames arrayFlags) {
	fmt.Println("Starting program successfully... Output will be displayed only when records are updated.")
	//Infinite loop to update records over time
	for {
		//GETS current record information
		url := "https://api.cloudflare.com/client/v4/zones/" + zone + "/dns_records"
		body := httpRequest(client, "GET", url, nil, email, token, zone)
		jsonData := unjsonify(body)

		cfRecords := map[string]Result
		for _, record := range jsonData.Result {
			if recordType == "A" {
				cfRecords[record.Name] = record
			}
		}

		publicIP := getIP()
		for _, name := range recordNames {
			jsonData := []byte{}
			if _, ok := cfRecords[name]; !ok {
				jsonData = jsonify("A", name, publicIP, false)
				httpRequest(client, "POST", url, jsonData, email, token, zone)
				fmt.Println("Added Record: " + name + " Updated IP: " + publicIP)
			} else {
				// Update record if it does not match current IP
				if cfRecords[name].Content != publicIP {
					jsonData = jsonify(cfRecords[name].Type, recordName, publicIP, cfRecords[name].Proxied)
					recordURL := url + "/" + cfRecords[name].Identifier
					httpRequest(client, "PUT", recordURL, jsonData, email, token, zone)
					fmt.Println("Updated Record: " + name + " Updated IP: " + publicIP)
				}
			}
		}
		time.Sleep(interval * time.Second) //Sleeping for n seconds
	}
}

func main() {
	getCredentials()

	flag.Var(&recordToUpdate, "r", "Records (sub.domain.tld) to update with current IP")
	flag.Parse()

	timeout := time.Duration(120 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}

	go update(client, recordsToUpdate)
}
