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

func TestCloudfrontInvalidCredentials(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudfront{}
	ch := make(chan *RequestError, 1)
	v = append(v, &Cloudfront{
		AccessKey:      "3f3f3f3rf",
		SecretKey:      "46y4g5454",
		DistributionID: os.Getenv("AWS_CLOUDFRONT_ID"),
		Resources: []string{
			"/index.xml",
		},
	})
	invalidateCloudfront(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 1 {
		t.Fatalf("Should contain 1 error(s), current errors: %d", len(ch))
	}
}

func TestCloudfrontInvalidDistributionID(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudfront{}
	ch := make(chan *RequestError, 1)
	v = append(v, &Cloudfront{
		AccessKey:      os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:      os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DistributionID: "23f3f34f34f3f2d24f34f4f34f344f43f",
		Resources: []string{
			"/index.xml",
		},
	})
	invalidateCloudfront(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 1 {
		t.Fatalf("Should contain 1 error(s), current errors: %d", len(ch))
	}
}

func TestCloudfrontWithoutResources(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	v := []*Cloudfront{}
	ch := make(chan *RequestError, 1)
	v = append(v, &Cloudfront{
		AccessKey:      os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretKey:      os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DistributionID: os.Getenv("AWS_CLOUDFRONT_ID"),
	})
	invalidateCloudfront(v, &wg, ch)
	wg.Wait()
	close(ch)
	if len(ch) != 1 {
		t.Fatalf("Should contain 1 error(s), current errors: %d", len(ch))
	}
}
