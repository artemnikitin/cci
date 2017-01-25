package main

import "bytes"

const (
	AWS = "AWS Cloudfront"
	CF  = "Cloudflare"
)

// RequestError represents container for error details
type RequestError struct {
	ID      string
	service string
	message string
}

func (e *RequestError) toString() string {
	var b bytes.Buffer
	b.WriteString("Error in service: ")
	b.WriteString(e.service)
	b.WriteString(", with ID: ")
	b.WriteString(e.ID)
	b.WriteString(", because of: ")
	b.WriteString(e.message)
	return b.String()
}

func newError(typ, id, message string) *RequestError {
	return &RequestError{
		ID:      id,
		service: typ,
		message: message,
	}
}
