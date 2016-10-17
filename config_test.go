package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConfigFromUnexistedFile(t *testing.T) {
	cases := []string{
		"",
		"-1",
		"/terer/dfdf/c.json",
		"file://dfdfd/v/v/c.json",
		"s3.amazonaws.com/artnik-test/test3.json",
	}

	for _, v := range cases {
		t.Run("Path: "+v, func(t *testing.T) {
			_, err := getConfig(v)
			if err == nil {
				t.Error("Should return error for trying to open config in unexisted file:", v)
			}
		})
	}
}

func TestGetConfigFileNotJSON(t *testing.T) {
	_, err := getConfig("/test1.json")
	if err == nil {
		t.Error("Should return error for trying to open config in non-JSON file")
	}
}

func TestGetConfigIncorrectJSON(t *testing.T) {
	_, err := getConfig("/test2.json")
	if err == nil {
		t.Error("Should return error for trying to open config with incorrct JSON structure")
	}
}

func TestGetConfigCorrect(t *testing.T) {
	server := serveJSON()
	defer server.Close()

	input := []string{"test3.json", server.URL}
	for _, v := range input {
		t.Run("Path: "+v, func(t *testing.T) {
			config, err := getConfig(v)
			if err != nil {
				t.Error("Error:", err, "on open config:", v)
			}
			if len(config.Cloudfront) != 2 {
				t.Error(config.Cloudfront)
				t.Error("Config should contain 2 entry for Cloudfront")
			}
			if len(config.Cloudflare) != 2 {
				t.Error(config.Cloudflare)
				t.Error("Config should contain 2 entry for Cloudflare")
			}
		})
	}
}

func serveJSON() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		bytes, _ := ioutil.ReadFile("test3.json")
		fmt.Fprintln(w, string(bytes))
	}))
	return server
}