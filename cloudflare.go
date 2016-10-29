package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

func invalidateCloudflare(data []*Cloudflare, wg *sync.WaitGroup, ch chan<- *RequestError) {
	for _, v := range data {
		go func(v *Cloudflare, wg *sync.WaitGroup) {
			fmt.Println("Starting Cloudflare cache invalidation for", v.ZoneID)
			defer wg.Done()

			opt := cloudflare.HTTPClient(&http.Client{
				Timeout: 10 * time.Second,
			})
			api, err := cloudflare.New(v.Key, v.Email, opt)
			if err != nil {
				ch <- newError(CF, v.ZoneID, err.Error())
				return
			}

			if v.PurgeAll == true {
				_, err := api.PurgeEverything(v.ZoneID)
				if err != nil {
					ch <- newError(CF, v.ZoneID, err.Error())
				}
			} else {
				req := cloudflare.PurgeCacheRequest{}
				for _, v := range v.Resources {
					req.Files = append(req.Files, v)
				}
				_, err := api.PurgeCache(v.ZoneID, req)
				if err != nil {
					ch <- newError(CF, v.ZoneID, err.Error())
				}
			}
		}(v, wg)
	}
}
