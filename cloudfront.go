package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

var (
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
)

func invalidateCloudfront(data []*Cloudfront, wg sync.WaitGroup) {
	for _, v := range data {
		go func(v *Cloudfront, wg sync.WaitGroup) {
			fmt.Println("Starting Cloudfront cache invalidation for ", v.DistributionID)
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
			s, err := session.NewSession(config)
			if err != nil {
				log.Println("AWS Error:", err)
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
			if err != nil {
				log.Println("Receive an error:", err, "on Cloudfront with distribution id:", v.DistributionID)
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
