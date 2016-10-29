package main

import "bytes"

type RequestError struct {
	ID      string
	service string
	message string
}

func (e *RequestError) ToString() string {
	var b bytes.Buffer
	b.WriteString("Error in service: ")
	b.WriteString(e.service)
	b.WriteString(" with ID: ")
	b.WriteString(e.ID)
	b.WriteString(" because: ")
	b.WriteString(e.message)
	return b.String()
}

func NewAWSCloudfrontError(id, message string) *RequestError {
	return &RequestError{
		ID:      id,
		service: "AWS Cloudfront",
		message: message,
	}
}

func NewCloudflareError(id, message string) *RequestError {
	return &RequestError{
		ID:      id,
		service: "Cloudflare",
		message: message,
	}
}
