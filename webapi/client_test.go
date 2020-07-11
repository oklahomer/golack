package webapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

const (
	AuthHeaderName   = "Authorization"
	AuthBearerSchema = "Bearer "
)

type GetResponseDummy struct {
	APIResponse
	Foo string
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()
	if config.RequestTimeout == 0 {
		t.Error("Default timeout is not set.")
	}
}

func TestWithHTTPClient(t *testing.T) {
	httpClient := &http.Client{}
	option := WithHTTPClient(httpClient)
	client := &Client{}

	option(client)

	if client.httpClient != httpClient {
		t.Errorf("Specified htt.Client is not set")
	}
}

func TestNewClient(t *testing.T) {
	config := &Config{Token: "abc", RequestTimeout: 1 * time.Second}
	optionCalled := false
	client := NewClient(config, func(*Client) { optionCalled = true })

	if client == nil {
		t.Fatal("Returned client is nil.")
	}

	if client.config != config {
		t.Errorf("Returned client does not have assigned config: %#v.", client.config)
	}

	if !optionCalled {
		t.Error("ClientOption is not called.")
	}

	if client.httpClient != http.DefaultClient {
		t.Errorf("When WithHTTPClient is not given, http.DefaultClient must be set: %+v", client.httpClient)
	}
}

func Test_buildEndpoint(t *testing.T) {
	params := url.Values{
		"foo": []string{"bar", "buzz"},
	}

	method := "rtm.start"
	endpoint := buildEndpoint(method, params)

	if endpoint == nil {
		t.Fatal("url is not returned.")
	}

	fooParam, _ := endpoint.Query()["foo"]
	if fooParam == nil || fooParam[0] != "bar" || fooParam[1] != "buzz" {
		t.Errorf("expected query parameter was not returned: %#v.", fooParam)
	}
}

func TestClient_Get(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := "abc"
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get(AuthHeaderName)
			if len(auth) == 0 {
				t.Fatal("Authorization header is not given")
			}

			tokenVal := auth[len(AuthBearerSchema):]
			if tokenVal != token {
				t.Errorf("Expected token value is not given: %s", auth)
			}

			if req.URL.Query().Get("bar") != "buzz" {
				t.Errorf("Expected query parameter is not given: %+v", req.URL.Query())
			}

			w.WriteHeader(http.StatusOK)

			response := &GetResponseDummy{
				APIResponse: APIResponse{OK: true},
				Foo:         "bar",
			}
			bytes, _ := json.Marshal(response)
			w.Write(bytes)
		})
		client := &Client{
			config: &Config{
				Token:          token,
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: &localRoundTripper{mux: mux}},
		}

		queryParams := url.Values{}
		queryParams.Set("bar", "buzz")
		returnedResponse := &GetResponseDummy{}
		err := client.Get(context.TODO(), "foo", queryParams, returnedResponse)

		if err != nil {
			t.Errorf("something went wrong. %#v", err)
		}

		if returnedResponse.OK != true {
			t.Errorf("OK status is wrong. %#v", returnedResponse)
		}

		if returnedResponse.Foo != "bar" {
			t.Errorf("foo value is wrong. %#v", returnedResponse)
		}
	})

	t.Run("status error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		err := client.Get(context.TODO(), "foo", nil, &GetResponseDummy{})

		if err == nil {
			t.Error("error should return when 404 is given.")
		}
	})

	t.Run("JSON error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		err := client.Get(context.TODO(), "foo", nil, &GetResponseDummy{})

		if err == nil {
			t.Error("error should return")
		}
	})
}

func TestClient_Post(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := "abc"
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get(AuthHeaderName)
			if len(auth) == 0 {
				t.Fatal("Authorization header is not given")
			}

			tokenVal := auth[len(AuthBearerSchema):]
			if tokenVal != token {
				t.Errorf("Expected token value is not given: %s", auth)
			}

			defer req.Body.Close()
			bytes, _ := ioutil.ReadAll(req.Body)
			query, _ := url.ParseQuery(string(bytes))
			if query.Get("foo") != "bar" {
				t.Errorf("Expected parameter is not passed: %+v", query)
			}

			w.WriteHeader(http.StatusOK)

			response := &APIResponse{OK: true}
			bytes, _ = json.Marshal(response)
			w.Write(bytes)
		})

		client := &Client{
			config: &Config{
				Token:          token,
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: &localRoundTripper{mux: mux}},
		}

		returnedResponse := &APIResponse{}
		values := url.Values{}
		values.Set("foo", "bar")
		err := client.Post(context.TODO(), "foo", values, returnedResponse)

		if err != nil {
			t.Errorf("something is wrong. %#v", err)
		}
	})

	t.Run("status error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		returnedResponse := &APIResponse{}
		err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

		if err == nil {
			t.Error("error should return when 500 is given.")
		}
	})

	t.Run("JSON error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		returnedResponse := &APIResponse{}
		err := client.Post(context.TODO(), "foo", url.Values{}, returnedResponse)

		if err == nil {
			t.Error("error should return")
		}
	})
}

func TestClient_PostJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token := "abc"
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			auth := req.Header.Get(AuthHeaderName)
			if len(auth) == 0 {
				t.Fatal("Authorization header is not given")
			}

			tokenVal := auth[len(AuthBearerSchema):]
			if tokenVal != token {
				t.Errorf("Expected token value is not given: %s", auth)
			}

			defer req.Body.Close()
			message := &PostMessage{}
			json.NewDecoder(req.Body).Decode(message)
			if message.Text != "hello" {
				t.Errorf("Expected parameter is not passed: %+v", message)
			}

			w.WriteHeader(http.StatusOK)

			response := &APIResponse{OK: true}
			bytes, _ := json.Marshal(response)
			w.Write(bytes)
		})

		client := &Client{
			config: &Config{
				Token:          token,
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: &localRoundTripper{mux: mux}},
		}

		returnedResponse := &APIResponse{}
		payload := &PostMessage{Text: "hello"}
		err := client.PostJSON(context.TODO(), "foo", payload, returnedResponse)

		if err != nil {
			t.Errorf("something is wrong. %#v", err)
		}
	})

	t.Run("status error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		returnedResponse := &APIResponse{}
		err := client.PostJSON(context.TODO(), "foo", &PostMessage{}, returnedResponse)

		if err == nil {
			t.Error("error should return when 500 is given.")
		}
	})

	t.Run("JSON error", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/api/foo", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid json"))
		})
		tripper := &localRoundTripper{mux: mux}
		client := &Client{
			config: &Config{
				Token:          "abc",
				RequestTimeout: 3 * time.Second,
			},
			httpClient: &http.Client{Transport: tripper},
		}

		returnedResponse := &APIResponse{}
		err := client.PostJSON(context.TODO(), "foo", &PostMessage{}, returnedResponse)

		if err == nil {
			t.Error("error should return")
		}
	})
}

type localRoundTripper struct {
	mux *http.ServeMux
}

func (l *localRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	l.mux.ServeHTTP(w, req)
	return w.Result(), nil
}
