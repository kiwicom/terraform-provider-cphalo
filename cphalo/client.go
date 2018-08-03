package cphalo

import (
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout = 180 * time.Second
	DefaultBaseUrl = "https://api.cloudpassage.com/v1"
)

type Client struct {
	AppKey    string
	AppSecret string
	BaseUrl   *url.URL
	client    *http.Client
	Timeout   time.Duration
}

func newClient(appKey string, appSecret string) *Client {
	baseUrl, _ := url.Parse(DefaultBaseUrl)
	c := &Client{AppKey: appKey, AppSecret: appSecret, BaseUrl: baseUrl}
	c.client = http.DefaultClient

	return c
}
