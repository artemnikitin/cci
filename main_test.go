package main

import (
	"log"
	"testing"
)

func TestEnd2End(t *testing.T) {
	config, err := getConfig("test.json")
	if err != nil {
		t.Fatal("Open config error: ", err.Error())
	}
	errors := make(chan *RequestError, config.getSize())
	execute(config, errors)
	if len(errors) != 0 {
		for v := range errors {
			log.Println(v.toString())
		}
		t.Error("Expected 0 errors, actual: ", len(errors))
	}
}
