# cci
[![Go Report Card](https://goreportcard.com/badge/github.com/artemnikitin/cci)](https://goreportcard.com/report/github.com/artemnikitin/cci)   [![codebeat badge](https://codebeat.co/badges/41e5be9d-a578-4bd0-87d4-5ffc564e38f0)](https://codebeat.co/projects/github-com-artemnikitin-cci)   [![Build Status](https://travis-ci.org/artemnikitin/cci.svg?branch=master)](https://travis-ci.org/artemnikitin/cci)   
CDN Cache Invalidation Tool

#### Description
Tool for invalidating cache for several CDN providers. Currently supports:
```
Cloudfront (AWS)
Cloudflare
```

#### Get it 
``` 
go get github.com/artemnikitin/cci
``` 

#### Use it
```
cci -config /path/to/config
cci -config https://example.com/config.json
```
Parameters:
- ``config`` specified path to config on hard drive or URL
- ``debug`` print additional info for debug, optional

#### Config 
Should be present as JSON file.
```json
{
	"cloudfront" : [{
		"access_key": "AWS access key",
		"secret_key": "AWS secret key",
		"distribution_id": "Cloudfront distribution ID",
		"resources": [
			"List of files for invalidation, optional",
			"Format: /index.html or /folder/*"
		]
	},
	{
		"access_key": "AWS access key",
		"secret_key": "AWS secret key",
		"region": "AWS region, format: eu-west-1",
		"distribution_id": "Cloudfront distribution ID"
	}],
	"cloudflare" : [{
		"email": "Cloudflare email",
		"key": "Cloudflare API key",
		"zone_id": "Cloudflare Zone ID",
		"purge_all": true
	},
	{
		"email": "Cloudflare email",
		"key": "Cloudflare API key",
		"zone_id": "Cloudflare Zone ID",
		"resources": [
			"List of files for invalidation, optional",
			"Format: http://example.com/index.html"
		]
	}]
}
```
