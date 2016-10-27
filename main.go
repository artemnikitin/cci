package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var (
	path  = flag.String("config", "", "Path to config in form of file or URL, required")
	debug = flag.Bool("debug", false, "Print debug info, optional")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().UnixNano())

	if *path == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conf, err := getConfig(*path)
	if err != nil {
		log.Fatal("Can't process config: ", err)
	}

	var wg sync.WaitGroup
	if len(conf.Cloudflare) > 0 {
		wg.Add(len(conf.Cloudflare))
		invalidateCloudflare(conf.Cloudflare, &wg)
	}
	if len(conf.Cloudfront) > 0 {
		wg.Add(len(conf.Cloudfront))
		invalidateCloudfront(conf.Cloudfront, &wg)
	}
	wg.Wait()
	fmt.Println("Done!")
}
