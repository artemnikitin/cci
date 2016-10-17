package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/cloudflare/cloudflare-go"
)

func invalidateCloudflare(data []*Cloudflare, wg sync.WaitGroup) {
	for _, v := range data {
		go func(v *Cloudflare, wg sync.WaitGroup) {
			fmt.Println("Starting Cloudflare cache invalidation for ", v.ZoneID)
			defer wg.Done()
			api, err := cloudflare.New(v.Key, v.Email)
			if err != nil {
				log.Println("Get an error on autentification in Cloudflare", err)
				return
			}
			if v.PurgeAll == true {
				resp, err := api.PurgeEverything(v.ZoneID)
				if err != nil {
					log.Println("Receive an error:", err.Error(), "on Cloudflare for zone:", v.ZoneID)
				}
				if *debug {
					log.Println(resp)
				}
			} else {
				if v.Resources != nil {
					req := cloudflare.PurgeCacheRequest{}
					for _, v := range v.Resources {
						req.Files = append(req.Files, v)
					}
					resp, err := api.PurgeCache(v.ZoneID, req)
					if err != nil {
						log.Println("Receive an error:", err.Error(), "on Cloudflare for zone:", v.ZoneID)
					}
					if *debug {
						log.Println(resp)
					}
				}
			}
		}(v, wg)
	}
}
