package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/config"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/publicip"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/secret"
)

func main() {
	defaultConfigDir := "/etc/cloudflare-dynamic-dns"
	cfgPath := flag.String("config", fmt.Sprintf("%s/config.json", defaultConfigDir), "Path to config file")
	scrtPath := flag.String("secret", fmt.Sprintf("%s/secret.json", defaultConfigDir), "Path to secret file")
	flag.Parse()

	scrtFile, err := os.ReadFile(*scrtPath)
	if err != nil {
		log.Fatalf("failed to read secret file: %s\n", err)
	}

	scrt, err := secret.NewSecret(scrtFile)
	if err != nil {
		log.Fatalf("failed to parse secret: %s\n", err)
	}

	cf, err := cloudflare.NewClient(scrt)
	if err != nil {
		log.Fatalf("failed to create Cloudflare client: %s\n", err)
	}

	cfgFile, err := os.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("failed to read config file: %s\n", err)
	}

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		log.Fatalf("failed to parse config: %s\n", err)
	}

	log.Printf("ensuring %d records", len(cfg.Records))

	cfgZones := cfg.GetZones()
	cfgRecords := cfg.GetRecords()
	previousIp := "startup"
	for {
		publicIp, err := publicip.Lookup()
		if err != nil {
			log.Printf("failed to get public ip: %s\n", err)
		}

		if previousIp != publicIp {
			log.Printf("ip change detected: %s -> %s\n", previousIp, publicIp)
			cfRecords, err := cf.GetDnsRecords(cfgZones)
			if err != nil {
				log.Printf("failed to get cloudflare dns records: %s\n", err)
			}

			rMap := map[string]cloudflare.CfRecord{}
			for _, record := range cfRecords {
				rMap[record.Name] = record
			}

			for _, r := range cfgRecords {
				if _, ok := rMap[r.Domain]; !ok {
					if err := cf.CreateRecord(r, publicIp); err != nil {
						log.Printf("failed to create record: %s\n", err)
					}
					log.Printf("added record: %s\n", r.Domain)
				} else {
					// Prevent additional calls to cloudflare if the IP hasn't changed, for example, on application startup
					if rMap[r.Domain].Content != publicIp {
						if err := cf.UpdateRecord(r, rMap[r.Domain].Id, publicIp); err != nil {
							log.Printf("failed to update record: %s\n", err)
						}
						log.Printf("updated record: %s\n", r.Domain)
					}
				}
			}
			previousIp = publicIp
		}
		time.Sleep(120 * time.Second)
	}
}
