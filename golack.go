// Package golack provides higher-level interfaces to interact with Events API, Web API and RTM API.
// Since its aim is to provides higher-level interfaces for easier setup, this lacks customizability.
// Directly use below sub-packages when customization is required.
//   - eventsapi ... Events API ( https://api.slack.com/events-api )
//   - webapi ... Web API ( https://api.slack.com/web )
//   - rtmapi ... RTM API ( https://api.slack.com/rtm )
package golack

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/eventsapi"
	"golang.org/x/xerrors"
	"net/http"
	"net/url"
	"time"

	"github.com/oklahomer/golack/rtmapi"
	"github.com/oklahomer/golack/webapi"
)

// Config defines the serializable/deserializable configuration value.
// Use NewConfig() to build a Config with default setting, and then pass the instance to json.Unmarshal()/yaml.Unmarshal() to apply customization.
type Config struct {
	AppSecret      string        `json:"app_secret" yaml:"app_secret"`
	Token          string        `json:"token" yaml:"token"`
	ListenPort     int           `json:"listen_port" yaml:"listen_port"`
	RequestTimeout time.Duration `json:"request_timeout" yaml:"request_timeout"`
}

// NewConfig returns initialized Config struct with default settings.
// AppSecret and Token are empty at this point. They can be set/updated by feeding this instance to json.Unmarshal/yaml.Unmarshal
// or by direct assignment.
func NewConfig() *Config {
	return &Config{
		AppSecret:      "",
		Token:          "",
		ListenPort:     8080,
		RequestTimeout: 3 * time.Second,
	}
}

// WebClient defines an interface that abstracts the webapi.Client.
type WebClient interface {
	Get(ctx context.Context, slackMethod string, queryParams url.Values, response interface{}) error
	Post(ctx context.Context, slackMethod string, payload url.Values, response interface{}) error
	PostJSON(ctx context.Context, slackMethod string, payload interface{}, response interface{}) error
}

type Option func(*Golack)

// WithWebClient provides a way to use pre-configured WebClient implementation.
// Pass the returned Option to New().
func WithWebClient(wc WebClient) Option {
	return func(g *Golack) {
		g.WebClient = wc
	}
}

// Golack works as a kind of facade to provide higher level interface to work with Events API, Web API and RTM API.
// For more customizability, use each sub-package that corresponds to each API.
type Golack struct {
	WebClient WebClient
	config    *Config
}

// New builds a new Golack instance with given config and options.
func New(config *Config, options ...Option) *Golack {
	g := &Golack{}
	g.config = config

	// Apply options to change specific behaviors
	for _, opt := range options {
		opt(g)
	}

	// If WebClient is not set with Option, then built one with default settings
	if g.WebClient == nil {
		apiConfig := webapi.NewConfig()
		apiConfig.Token = g.config.Token
		apiConfig.RequestTimeout = g.config.RequestTimeout
		g.WebClient = webapi.NewClient(apiConfig)
	}

	return g
}

// PostMessage posts a postMessage to Slack.
//
// See https://api.slack.com/methods/chat.postMessage for official document.
func (g *Golack) PostMessage(ctx context.Context, postMessage *webapi.PostMessage) (*webapi.APIResponse, error) {
	response := &webapi.APIResponse{}
	err := g.WebClient.PostJSON(ctx, "chat.postMessage", postMessage, response)
	if err != nil {
		return nil, err
	}

	if response.OK != true {
		return nil, fmt.Errorf("failed chat.postMessage request: %s", response.Error)
	}

	return response, nil
}

// ConnectRTM connects to Slack WebSocket server.
//
// See https://api.slack.com/rtm for official document.
func (g *Golack) ConnectRTM(ctx context.Context) (rtmapi.Connection, error) {
	rtmStart := &webapi.RTMStart{}
	err := g.WebClient.Get(ctx, "rtm.start", nil, rtmStart)
	if err != nil {
		return nil, err
	}

	if rtmStart.OK != true {
		return nil, fmt.Errorf("failed rtm.start request: %s", rtmStart.Error)
	}

	return rtmapi.Connect(ctx, rtmStart.URL)
}

// RunServer starts a server to interact with Events API.
// The server runs in another goroutine so this method is not blocking.
// To pass and notify the error state of the server from the server, this returns a channel that passes the error.
// When the error is returned from the channel, the server is not running or is already stopped.
//
// See https://api.slack.com/events-api for official document.
func (g *Golack) RunServer(ctx context.Context, receiver eventsapi.EventReceiver) <-chan error {
	errChan := make(chan error, 1)

	// Setup a request validator
	// For better security, this checks each request's signature
	appSecret := g.config.AppSecret
	if appSecret == "" {
		errChan <- xerrors.New("application secret is not set")
		return errChan
	}
	optValidator := eventsapi.WithRequestValidator(&eventsapi.SignatureValidator{Secret: appSecret})

	// Setup server and run it
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", g.config.ListenPort),
		Handler: eventsapi.SetupHandler(receiver, optValidator),
	}
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	// Shutdown the server
	go func() {
		<-ctx.Done()
		//noinspection ALL
		srv.Shutdown(ctx)
	}()

	return errChan
}
