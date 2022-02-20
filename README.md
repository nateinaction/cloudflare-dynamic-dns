# cloudflare-dynamic-dns
[![ci](https://github.com/nateinaction/cloudflare-dynamic-dns/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/nateinaction/cloudflare-dynamic-dns/actions/workflows/ci.yaml)

`cloudflare-dynamic-dns` is a dynamic dns client that can be used to update DNS records in Cloudflare.

## Use

```
$ cloudflare-dynamic-dns -r <domain.com> -r <sub1.domain.com> -r <sub2.domain.com>
```
### Environment Variables
The following environment variables are required:

| Variable | Function |
| ---- | ---- | 
| CF_EMAIL | The email address associated with your Cloudflare account | 
| CF_TOKEN | A [Cloudflare API token](https://developers.cloudflare.com/api/tokens/create) | 
| CF_ZONE | The zone ID for your domain | 
