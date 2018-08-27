package cphalo

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout    = 180 * time.Second
	DefaultBaseUrl    = "https://api.cloudpassage.com"
	DefaultApiVersion = "v1"
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

func (pc *Client) Validate() (bool, error) {
	return true, nil
}

func (pc *Client) AuthRequest(appKey string, appSecret string) (*http.Request, error) {
	rsc := "/oauth/access_token?grant_type=client_credentials"
	method := "POST"
	baseUrl, err := url.Parse(pc.BaseUrl.String() + rsc)
	log.Println("Going to authenticate and obtain access token.")
	log.Println("Auth URL is " + baseUrl.String())
	if err != nil {
		return nil, err
	}

	authString := appKey + ":" + appSecret
	encodedAuthString := b64.StdEncoding.EncodeToString([]byte(authString))

	req, err := http.NewRequest(method, baseUrl.String(), nil)
	req.Header.Add("Authorization", "Basic "+encodedAuthString)

	return req, err
}

func (pc *Client) NewRequest(method string, rsc string, params map[string]string) (*http.Request, error) {
	baseUrl, err := url.Parse(pc.BaseUrl.String() + "/" + DefaultApiVersion + "/" + rsc)
	if err != nil {
		return nil, err
	}

	if params != nil {
		ps := url.Values{}
		for k, v := range params {
			ps.Set(k, v)
		}
		baseUrl.RawQuery = ps.Encode()
	}

	req, err := pc.AuthRequest(pc.AppKey, pc.AppSecret)
	if err != nil {
		return nil, err
	}

	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	m := &apiKeyJsonResponse{}
	err = json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return nil, err
	}

	log.Println("Going to make a new request.")
	log.Println("Base URL is " + baseUrl.String())

	req, err = http.NewRequest(method, baseUrl.String(), nil)
	req.Header.Add("Authorization", "Bearer "+m.AccessToken)

	return req, err
}

func (pc *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := pc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return resp, err
	}

	err = parseResponse(resp, v)

	return resp, err
}

func parseResponse(r *http.Response, v interface{}) error {
	if v == nil {
		return fmt.Errorf("nil interface provided to decodeResponse")
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	err := json.Unmarshal([]byte(bodyString), &v)

	return err
}

func validateResponse(r *http.Response) error {
	log.Println("Reponse code is " + r.Status)
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	log.Println("Wrong processing")
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	m := &errorJsonResponse{}
	err := json.Unmarshal([]byte(bodyString), &m)
	if err != nil {
		return err
	}

	return m.Error
}
