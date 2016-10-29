package main

import (
	"os"
	"sync"
	"testing"
)

func TestCloudflareInvalidCredentials(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudflare{}
	ch := make(chan *RequestError, 1)
	v = append(v, &Cloudflare{
		Email:  "3g34g3434",
		Key:    "wf34f3f34f4",
		ZoneID: os.Getenv("CLOUDFLARE_ID"),
		Resources: []string{
			"https://artemnikitin.com/index.xml",
		},
	})
	invalidateCloudflare(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 1 {
		t.Fatalf("Should contain 1 error(s), current errors: %d", len(ch))
	}
}

func TestCloudflareInvalidZone(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	v := []*Cloudflare{}
	ch := make(chan *RequestError, 2)
	v = append(v, &Cloudflare{
		Email:  os.Getenv("CLOUDFLARE_EMAIL"),
		Key:    os.Getenv("CLOUDFLARE_KEY"),
		ZoneID: "3t35g5gg42ed232d",
		Resources: []string{
			"https://artemnikitin.com/index.xml",
		},
	}, &Cloudflare{
		Email:    os.Getenv("CLOUDFLARE_EMAIL"),
		Key:      os.Getenv("CLOUDFLARE_KEY"),
		ZoneID:   "3t35g5gg42ed232d",
		PurgeAll: true,
	})
	invalidateCloudflare(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 2 {
		t.Fatalf("Should contain 2 error(s), current errors: %d", len(ch))
	}
}
