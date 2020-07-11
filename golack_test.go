package golack

import (
	"context"
	"fmt"
	"github.com/oklahomer/golack/eventsapi"
	"github.com/oklahomer/golack/testutil"
	"github.com/oklahomer/golack/webapi"
	"golang.org/x/xerrors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type DummyWebClient struct {
	GetFunc      func(ctx context.Context, slackMethod string, queryParams url.Values, response interface{}) error
	PostFunc     func(ctx context.Context, slackMethod string, payload url.Values, response interface{}) error
	PostJSONFunc func(ctx context.Context, slackMethod string, payload interface{}, response interface{}) error
}

func (wc *DummyWebClient) Get(ctx context.Context, slackMethod string, queryParams url.Values, response interface{}) error {
	return wc.GetFunc(ctx, slackMethod, queryParams, response)
}

func (wc *DummyWebClient) Post(ctx context.Context, slackMethod string, payload url.Values, response interface{}) error {
	return wc.PostFunc(ctx, slackMethod, payload, response)
}

func (wc *DummyWebClient) PostJSON(ctx context.Context, slackMethod string, payload interface{}, response interface{}) error {
	return wc.PostJSONFunc(ctx, slackMethod, payload, response)
}

type DummyReceiver struct {
	ReceiveFunc func(wrapper *eventsapi.EventWrapper)
}

func (d DummyReceiver) Receive(wrapper *eventsapi.EventWrapper) {
	d.ReceiveFunc(wrapper)
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config == nil {
		t.Fatal("Returned *Config is nil.")
	}

	if config.RequestTimeout == 0 {
		t.Error("Default timeout is not set.")
	}

	if config.ListenPort == 0 {
		t.Error("Default listen port is not set.")
	}
}

func TestWithWebClient(t *testing.T) {
	webClient := &DummyWebClient{}
	option := WithWebClient(webClient)
	g := &Golack{}

	option(g)

	if g.WebClient != webClient {
		t.Errorf("Specified WebClient is not set.")
	}
}

func TestNew(t *testing.T) {
	config := &Config{}
	optionCalled := false

	g := New(config, func(_ *Golack) { optionCalled = true })

	if g == nil {
		t.Fatal("Returned *Golack is nil.")
	}

	if !optionCalled {
		t.Error("Option is not called.")
	}
}

func TestGolack_PostMessage(t *testing.T) {
	t.Run("Web API returns error status", func(t *testing.T) {
		expectedErr := xerrors.New("DUMMY")
		webClient := &DummyWebClient{
			PostJSONFunc: func(_ context.Context, _ string, _ interface{}, _ interface{}) error {
				return expectedErr
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.PostMessage(context.TODO(), &webapi.PostMessage{})
		if err == nil {
			t.Fatal("Error is not returned.")
		}
		if err != expectedErr {
			t.Fatalf("Expected error is not returned: %+v", err)
		}
	})

	t.Run("Web API returns error response", func(t *testing.T) {
		webClient := &DummyWebClient{
			PostJSONFunc: func(_ context.Context, _ string, _ interface{}, response interface{}) error {
				resp := response.(*webapi.APIResponse)
				resp.OK = false
				resp.Error = "some error"
				return nil
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.PostMessage(context.TODO(), &webapi.PostMessage{})
		if err == nil {
			t.Fatal("Expected error is not returned.")
		}
	})

	t.Run("success", func(t *testing.T) {
		webClient := &DummyWebClient{
			PostJSONFunc: func(_ context.Context, _ string, payload interface{}, response interface{}) error {
				message := payload.(*webapi.PostMessage)
				if message.ChannelID != "channel" && message.Text != "my message" {
					t.Errorf("Expected payload is not passed: %+v", message)
				}

				resp := response.(*webapi.APIResponse)
				resp.OK = true
				resp.Error = ""
				return nil
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		postMessage := webapi.NewPostMessage("channel", "my message")
		response, err := g.PostMessage(context.TODO(), postMessage)

		if err != nil {
			t.Errorf("something is wrong. %#v", err)
		}

		if response.OK != true {
			t.Errorf("OK status is wrong. %#v", response)
		}
	})
}

func TestGolack_ConnectRTM(t *testing.T) {
	t.Run("Web API returns error status", func(t *testing.T) {
		expectedErr := xerrors.New("DUMMY")
		webClient := &DummyWebClient{
			GetFunc: func(_ context.Context, _ string, _ url.Values, _ interface{}) error {
				return expectedErr
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.ConnectRTM(context.Background())
		if err == nil {
			t.Fatal("Error is not returned.")
		}
		if err != expectedErr {
			t.Fatalf("Expected error is not returned: %+v", err)
		}
	})

	t.Run("Web API returns error response", func(t *testing.T) {
		webClient := &DummyWebClient{
			GetFunc: func(_ context.Context, _ string, _ url.Values, response interface{}) error {
				resp := response.(*webapi.RTMStart)
				resp.OK = false
				resp.Error = "some error"
				return nil
			},
		}
		g := &Golack{
			WebClient: webClient,
		}

		_, err := g.ConnectRTM(context.Background())
		if err == nil {
			t.Fatal("Expected error is not returned.")
		}
	})

	t.Run("connect WebSocket server", func(t *testing.T) {
		testutil.RunWithWebSocket(func(addr net.Addr) {
			webClient := &DummyWebClient{
				GetFunc: func(_ context.Context, _ string, _ url.Values, response interface{}) error {
					resp := response.(*webapi.RTMStart)
					resp.OK = true
					resp.URL = fmt.Sprintf("ws://%s%s", addr, "/ping")
					resp.Error = ""
					return nil
				},
			}
			g := &Golack{
				WebClient: webClient,
			}

			rtm, err := g.ConnectRTM(context.Background())
			if err != nil {
				t.Fatalf("Unexpected error is returned: %s", err.Error())
			}

			err = rtm.Ping()
			if err != nil {
				t.Fatalf("Unexpected error is returned on Ping: %s", err.Error())
			}
		})
	})
}

func TestGolack_RunServer(t *testing.T) {
	t.Run("without app secret", func(t *testing.T) {
		g := &Golack{config: &Config{}}

		errCh := g.RunServer(context.Background(), &DummyReceiver{})

		select {
		case err := <-errCh:
			if !strings.Contains(err.Error(), "application secret") {
				t.Errorf("Unexpected error is returned: %s", err.Error())
			}

		case <-time.NewTimer(1 * time.Second).C:
			t.Fatal("Expected error is not returned.")
		}
	})

	t.Run("run", func(t *testing.T) {
		g := &Golack{
			config: &Config{
				AppSecret:      "DUMMY",
				Token:          "",
				ListenPort:     0, // Find a free port
				RequestTimeout: 0,
			},
		}

		ctx, cancel := context.WithCancel(context.Background())
		errCh := g.RunServer(ctx, &DummyReceiver{})

		select {
		case err := <-errCh:
			t.Fatalf("Unexpected error is returned: %s", err.Error())

		case <-time.NewTimer(100 * time.Millisecond).C:
			// O.K.
		}

		cancel()

		select {
		case err := <-errCh:
			if err != http.ErrServerClosed {
				t.Errorf("Expected type of error is not returned: %+v", err)
			}

		case <-time.NewTimer(100 * time.Millisecond).C:
			t.Error("Expected error is not returned.")
		}
	})
}
