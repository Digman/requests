package requests

import (
	"context"
	"time"

	tls_client "github.com/Digman/tls-client"
	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/websocket"
)

// WebSocket creates a WebSocket connection using the client's TLS fingerprint.
//
// Example:
//
//	conn, err := client.WebSocket("wss://example.com/ws")
//	if err != nil { ... }
//	defer conn.Close()
//
//	conn.WriteMessage(websocket.TextMessage, []byte("hello"))
//	_, msg, _ := conn.ReadMessage()
func (c *Client) WebSocket(wsUrl string, headers ...http.Header) (*websocket.Conn, error) {
	return c.WebSocketWithTimeout(wsUrl, 30000, headers...)
}

// WebSocketWithTimeout creates a WebSocket connection with a custom handshake timeout.
//
// timeout is in milliseconds.
func (c *Client) WebSocketWithTimeout(wsUrl string, timeout int, headers ...http.Header) (*websocket.Conn, error) {
	h := http.Header{}
	if len(headers) > 0 {
		h = headers[0]
	}

	h.Set("User-Agent", c.UserAgent)

	if len(c.HeaderOrder) > 0 {
		h[http.HeaderOrderKey] = c.HeaderOrder
	}

	opts := []tls_client.WebsocketOption{
		tls_client.WithTlsClient(c.tlsClient),
		tls_client.WithUrl(wsUrl),
		tls_client.WithHeaders(h),
		tls_client.WithHandshakeTimeoutMilliseconds(timeout),
	}

	if c.tlsClient.GetCookieJar() != nil {
		opts = append(opts, tls_client.WithCookiejar(c.tlsClient.GetCookieJar()))
	}

	ws, err := tls_client.NewWebsocket(nil, opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	return ws.Connect(ctx)
}

// WebSocketForceH1 creates a WebSocket connection with ForceHttp1 enabled.
// Required for servers that negotiate HTTP/2 via ALPN but need HTTP/1.1 for WebSocket.
func (c *Client) WebSocketForceH1(wsUrl string, timeout int, headers ...http.Header) (*websocket.Conn, error) {
	h := http.Header{}
	if len(headers) > 0 {
		h = headers[0]
	}

	h.Set("User-Agent", c.UserAgent)

	if len(c.HeaderOrder) > 0 {
		h[http.HeaderOrderKey] = c.HeaderOrder
	}

	// 創建獨立的 ForceHttp1 tls-client
	clientProfile := getClientProfile(c.UserAgent)
	h1Options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutMilliseconds(timeout),
		tls_client.WithClientProfile(clientProfile),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithForceHttp1(),
	}
	if c.RawProxy != "" {
		h1Options = append(h1Options, tls_client.WithProxyUrl(c.RawProxy))
	}
	h1Client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), h1Options...)
	if err != nil {
		return nil, err
	}

	opts := []tls_client.WebsocketOption{
		tls_client.WithTlsClient(h1Client),
		tls_client.WithUrl(wsUrl),
		tls_client.WithHeaders(h),
		tls_client.WithHandshakeTimeoutMilliseconds(timeout),
	}

	ws, err := tls_client.NewWebsocket(nil, opts...)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	defer cancel()

	return ws.Connect(ctx)
}
