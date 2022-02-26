package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/cloudflare"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/config"
	"github.com/nateinaction/cloudflare-dynamic-dns/pkg/network"
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

	cf := cloudflare.NewClient(scrt)

	cfgFile, err := os.ReadFile(*cfgPath)
	if err != nil {
		log.Fatalf("failed to read config file: %s\n", err)
	}

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		log.Fatalf("failed to parse config: %s\n", err)
	}

	previousIp := &network.Ip{
		V4: "startup",
	}
	for {
		publicIp, err := network.NewIp()
		if err != nil {
			log.Printf("failed to get public ip: %s\n", err)
		}

		cfgZones := cfg.GetZones()
		if !publicIp.Match(previousIp) {
			log.Printf("ip change detected: %s -> %s\n", previousIp.V4, publicIp.V4)
			cfRecords, err := cf.GetRecords(cfgZones)
			if err != nil {
				log.Printf("failed to get cloudflare dns records: %s\n", err)
			}

			cfgRecords := cfg.GetRecords(publicIp)
			for _, r := range cfgRecords {
				if _, ok := cfRecords[r.Name]; !ok {
					if err := cf.CreateRecord(r, cfgZones[r.ZoneId], publicIp); err != nil {
						log.Printf("failed to create record %s: %s\n", r.Name, err)
						continue
					}
					log.Printf("added record: %s\n", r.Name)
				} else {
					// Prevent additional calls to cloudflare if the record has not changed
					if !cfRecords[r.Name].Match(r) {
						record := cfRecords[r.Name]
						zone := cfgZones[cfRecords[r.Name].ZoneId]
						if err := cf.UpdateRecord(record, zone, publicIp); err != nil {
							log.Printf("failed to update record %s: %s\n", r.Name, err)
							continue
						}
						log.Printf("updated record: %s\n", r.Name)
					}
				}
			}
			previousIp = publicIp
		}
		time.Sleep(2 * time.Minute)
	}
}
