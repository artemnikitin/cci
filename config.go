package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Config represents config for tool
type Config struct {
	Cloudfront []*Cloudfront `json:"cloudfront"`
	Cloudflare []*Cloudflare `json:"cloudflare"`
}

// Cloudfront represents config details for specific service
type Cloudfront struct {
	AccessKey      string   `json:"access_key"`
	SecretKey      string   `json:"secret_key"`
	DistributionID string   `json:"distribution_id"`
	Resources      []string `json:"resources"`
	Region         string   `json:"region,omitempty"`
}

// Cloudflare represents config details for specific service
type Cloudflare struct {
	Email     string   `json:"email"`
	Key       string   `json:"key"`
	ZoneID    string   `json:"zone_id"`
	PurgeAll  bool     `json:"purge_all"`
	Resources []string `json:"resources,omitempty"`
}

func (c *Config) getSize() int {
	return len(c.Cloudflare) + len(c.Cloudfront)
}

func getConfig(path string) (*Config, error) {
	var conf *Config
	var err error
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		conf, err = getConfigFromURL(path)
	} else {
		conf, err = getConfigFromFile(path)
	}
	return conf, err
}

func getConfigFromURL(path string) (*Config, error) {
	conf := &Config{}
	resp, err := http.DefaultClient.Get(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func getConfigFromFile(path string) (*Config, error) {
	conf := &Config{}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
