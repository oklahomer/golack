package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

const (
	slackAPIEndpointFormat = "https://slack.com/api/%s"
)

type Config struct {
	Token          string        `json:"token" yaml:"token"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

func NewConfig() *Config {
	return &Config{
		Token:          "",
		RequestTimeout: 3 * time.Second,
	}
}

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config, options ...ClientOption) *Client {
	c := &Client{config: config}
	for _, opt := range options {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	return c
}

func buildEndpoint(slackMethod string, queryParams url.Values) *url.URL {
	requestURL, err := url.Parse(fmt.Sprintf(slackAPIEndpointFormat, slackMethod))
	if err != nil {
		panic(err.Error())
	}

	if queryParams != nil {
		requestURL.RawQuery = queryParams.Encode()
	}

	return requestURL
}

func (client *Client) Get(ctx context.Context, slackMethod string, queryParams url.Values, response interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, queryParams)
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.WithContext(reqCtx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Token))

	// Do request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	// Handle response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}

func statusErr(resp *http.Response) error {
	reqDump := []byte("N/A")
	if resp.Request != nil {
		dump, err := httputil.DumpRequestOut(resp.Request, true)
		if err != nil {
			reqDump = dump
		}
	}

	resDump, resErr := httputil.DumpResponse(resp, true)
	if resErr != nil {
		resDump = []byte("N/A")
	}

	return fmt.Errorf("response status error. Status: %d.\nRequest: %s\nResponse: %s", resp.StatusCode, string(reqDump), string(resDump))
}

func (client *Client) Post(ctx context.Context, slackMethod string, payload url.Values, response interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, nil)
	req, err := http.NewRequest("POST", endpoint.String(), strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.WithContext(reqCtx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Token))

	// Do request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	// Handle response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) PostJSON(ctx context.Context, slackMethod string, payload interface{}, response interface{}) error {
	// Prepare request
	endpoint := buildEndpoint(slackMethod, nil)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(body))
	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return err
	}
	reqCtx, cancel := context.WithTimeout(ctx, client.config.RequestTimeout)
	defer cancel()
	req.WithContext(reqCtx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.Token))

	// Do request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Usually, the API returns a JSON structure with status code 200.
	// https://api.slack.com/web#evaluating_responses
	if resp.StatusCode != http.StatusOK {
		return statusErr(resp)
	}

	// Handle response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}
	return nil
}
