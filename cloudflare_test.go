package main

import (
	"os"
	"sync"
	"testing"
)

func TestCloudflare(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudflare{}
	ch := make(chan *RequestError)
	v = append(v, &Cloudflare{
		Email:  os.Getenv("CLOUDFLARE_EMAIL"),
		Key:    os.Getenv("CLOUDFLARE_KEY"),
		ZoneID: os.Getenv("CLOUDFLARE_ID"),
		Resources: []string{
			"https://artemnikitin.com/index.xml",
		},
	})
	invalidateCloudflare(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 0 {
		t.Fatalf("Should contain 0 error(s), current errors: %d", len(ch))
	}
}
