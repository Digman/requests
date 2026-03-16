package requests

import (
	"context"

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
	return c.WebSocketWithTimeout(wsUrl, defaultTimeout, headers...)
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

	return ws.Connect(context.Background())
}
