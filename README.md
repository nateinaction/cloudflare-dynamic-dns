# cloudflare-dynamic-dns
[![ci](https://github.com/nateinaction/cloudflare-dynamic-dns/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/nateinaction/cloudflare-dynamic-dns/actions/workflows/ci.yaml)

`cloudflare-dynamic-dns` is a dynamic dns client that can be used to update DNS records in Cloudflare.

## Use

```
$ cloudflare-dynamic-dns -config <optional path to config file> -secret <optional path to secret file>
```
### Config File
On startup the service will look for a config file in `/etc/cloudflare-dynamic-dns/config.json`. The location of the config file can be overridden by specifying the `-config` flag. The config file should be a JSON file with the following format:

```json
{
	"records": [
		{
			"type": "A",
			"name": "example.com",
			"zone_id": "cloudflare_zone_id"
		},
		{
			"type": "A",
			"name": "sub.example.net",
			"zone_id": "cloudflare_zone_id",
			"proxy": true,
			"ttl": 300
		}
	]
}
```

### Secret File
On startup the service will look for a secret file in `/etc/cloudflare-dynamic-dns/secret.json`. The location of the secret file can be overridden by specifying the `-secret` flag. The secret file should be a JSON file with the following format:

```json
{
	"email": "cloudflare_account_email",
	"token": "cloudflare_api_token"
}
```