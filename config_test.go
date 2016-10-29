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

func TestGetConfigSize(t *testing.T) {
	cases := []struct {
		title  string
		size   int
		config *Config
	}{
		{
			"empty config",
			0,
			&Config{},
		},
		{
			"only 1 entity",
			1,
			&Config{
				Cloudfront: []*Cloudfront{{}},
			},
		},
		{
			"only 1 entity with several values",
			3,
			&Config{
				Cloudfront: []*Cloudfront{{}, {}, {}},
			},
		},
		{
			"every type of entity",
			2,
			&Config{
				Cloudfront: []*Cloudfront{{}},
				Cloudflare: []*Cloudflare{{}},
			},
		},
		{
			"every type of entity with several values",
			3,
			&Config{
				Cloudfront: []*Cloudfront{{}},
				Cloudflare: []*Cloudflare{{}, {}},
			},
		},
	}

	for _, v := range cases {
		t.Run(v.title, func(t *testing.T) {
			size := v.config.getSize()
			if size != v.size {
				t.Errorf("Expected size: %d, actual: %d", v.size, size)
			}
		})
	}
}

func TestGetConfigFileNotJSON(t *testing.T) {
	_, err := getConfig("/README.md")
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

func TestGetConfigIncorrect(t *testing.T) {
	server := serveJSON()
	defer server.Close()

	url := "http://example.com"
	t.Run("Path: "+url, func(t *testing.T) {
		_, err := getConfig(url)
		if err == nil {
			t.Error("Error:", err, "on open config:", url)
		}
	})
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
