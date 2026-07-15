package requests

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	tls_client "github.com/Digman/tls-client"
)

type proxyTestClient struct {
	tls_client.HttpClient
	setProxy func(string) error
	getProxy func() string
}

func (c *proxyTestClient) SetProxy(proxy string) error {
	return c.setProxy(proxy)
}

func (c *proxyTestClient) GetProxy() string {
	return c.getProxy()
}

func TestClientRequestUsesConfiguredContext(t *testing.T) {
	tests := []struct {
		name string
		send func(*Request) *Request
	}{
		{name: "form", send: func(r *Request) *Request { return r.Post("http://127.0.0.1:1").Send() }},
		{name: "json", send: func(r *Request) *Request {
			return r.SetJson(map[string]string{"x": "y"}).Post("http://127.0.0.1:1").Send()
		}},
		{name: "multipart", send: func(r *Request) *Request { return r.SetFileData("x", "y", false).Post("http://127.0.0.1:1").Send() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			client := NewClient("test")
			client.SetContext(ctx)
			_, _, err := tt.send(client.NewRequest()).End()
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("request error = %v, want context.Canceled", err)
			}
		})
	}
}

func TestClientContextCanChangeConcurrentlyWithNewRequest(t *testing.T) {
	client := NewClient("test")
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			client.SetContext(context.Background())
		}()
		go func() {
			defer wg.Done()
			if client.NewRequest().ctx == nil {
				t.Error("NewRequest() context = nil")
			}
		}()
	}
	wg.Wait()
}

func TestClientSetProxyPreservesRawInputAndFailureState(t *testing.T) {
	client := NewClient("test")
	if err := client.SetProxy("127.0.0.1:8080"); err != nil {
		t.Fatalf("SetProxy() error = %v", err)
	}
	rawProxy, proxyURL := client.RawProxy, client.ProxyUrl
	if rawProxy != "127.0.0.1:8080" || proxyURL == nil || proxyURL.Scheme != "http" {
		t.Fatalf("proxy metadata = %q, %#v", rawProxy, proxyURL)
	}
	oldRaw, oldURL := rawProxy, proxyURL.String()
	if err := client.SetProxy("http://[::1"); err == nil {
		t.Fatal("invalid SetProxy() error = nil")
	}
	rawProxy, proxyURL = client.RawProxy, client.ProxyUrl
	if rawProxy != oldRaw || proxyURL == nil || proxyURL.String() != oldURL {
		t.Fatalf("proxy metadata changed after failure: %q, %#v", rawProxy, proxyURL)
	}
}

func TestClientSetProxySerializesUnderlyingAndMetadata(t *testing.T) {
	client := NewClient("test")
	firstEntered := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondEntered := make(chan struct{})
	var underlyingMu sync.RWMutex
	var underlyingProxy string
	client.tlsClient = &proxyTestClient{
		HttpClient: client.tlsClient,
		setProxy: func(proxy string) error {
			underlyingMu.Lock()
			underlyingProxy = proxy
			underlyingMu.Unlock()
			switch proxy {
			case "http://127.0.0.1:8080":
				close(firstEntered)
				<-releaseFirst
			case "http://127.0.0.1:8081":
				close(secondEntered)
			}
			return nil
		},
		getProxy: func() string {
			underlyingMu.RLock()
			defer underlyingMu.RUnlock()
			return underlyingProxy
		},
	}

	firstDone := make(chan error, 1)
	go func() {
		firstDone <- client.SetProxy("127.0.0.1:8080")
	}()
	<-firstEntered

	secondDone := make(chan error, 1)
	secondCalling := make(chan struct{})
	go func() {
		close(secondCalling)
		secondDone <- client.SetProxy("127.0.0.1:8081")
	}()
	<-secondCalling
	secondStartedEarly := false
	select {
	case <-secondEntered:
		secondStartedEarly = true
	case <-time.After(50 * time.Millisecond):
	}
	close(releaseFirst)
	if err := <-firstDone; err != nil {
		t.Fatalf("first SetProxy() error = %v", err)
	}
	if err := <-secondDone; err != nil {
		t.Fatalf("second SetProxy() error = %v", err)
	}
	if secondStartedEarly {
		t.Fatal("second SetProxy entered the underlying client before the first metadata commit")
	}

	rawProxy, proxyURL := client.RawProxy, client.ProxyUrl
	if rawProxy != "127.0.0.1:8081" || proxyURL == nil || proxyURL.String() != "http://127.0.0.1:8081" {
		t.Fatalf("final proxy metadata = %q, %#v", rawProxy, proxyURL)
	}
	if got := client.tlsClient.GetProxy(); got != proxyURL.String() {
		t.Fatalf("underlying proxy = %q, metadata = %q", got, proxyURL.String())
	}
}
