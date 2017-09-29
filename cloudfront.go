package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func invalidateCloudfront(data []*Cloudfront, wg *sync.WaitGroup, ch chan<- *RequestError) {
	for _, v := range data {
		go func(v *Cloudfront, wg *sync.WaitGroup) {
			fmt.Println("Starting Cloudfront cache invalidation for", v.DistributionID)
			defer wg.Done()
			var list []*string
			if v.Resources != nil {
				for _, v := range v.Resources {
					list = append(list, aws.String(v))
				}
			}

			config := aws.NewConfig()
			config.WithCredentials(credentials.NewStaticCredentials(v.AccessKey, v.SecretKey, ""))
			if v.Region != "" {
				config.WithRegion(v.Region)
			}
			if *debug {
				config.WithLogLevel(aws.LogDebugWithHTTPBody)
			}
			config.WithHTTPClient(&http.Client{
				Timeout: 60 * time.Second,
			})

			s, err := session.NewSession(config)
			if err != nil {
				ch <- newError(AWS, v.DistributionID, err.Error())
				return
			}

			client := cloudfront.New(s)
			params := &cloudfront.CreateInvalidationInput{
				DistributionId: aws.String(v.DistributionID),
				InvalidationBatch: &cloudfront.InvalidationBatch{
					CallerReference: aws.String(randomString(12)),
					Paths: &cloudfront.Paths{
						Quantity: aws.Int64(int64(len(v.Resources))),
						Items:    list,
					},
				},
			}

			_, err = client.CreateInvalidation(params)
			if err != nil && !strings.Contains(err.Error(), "Client.Timeout exceeded") {
				ch <- newError(AWS, v.DistributionID, err.Error())
			}
		}(v, wg)
	}
}

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
