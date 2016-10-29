package main

import (
	"os"
	"sync"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	str1 := randomString(10)
	str2 := randomString(10)
	if str1 == str2 {
		t.Errorf("Generated strings: %s, %s should be different", str1, str2)
	}
}

func TestCloudfrontCorrect(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudfront{}
	ch := make(chan *RequestError)
	v = append(v, &Cloudfront{
		AccessKey:      os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:      os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DistributionID: os.Getenv("AWS_CLOUDFRONT_ID"),
		Resources: []string{
			"/index.xml",
		},
	})
	invalidateCloudfront(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 0 {
		t.Fatalf("Should contain 0 error(s), current errors: %d", len(ch))
	}
}
