package requests

import (
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"net/url"
	"strings"
)

type Client struct {
	tlsClient tls_client.HttpClient

	ProxyUrl   *url.URL
	RawProxy   string
	UserAgent  string
	WindowSize [2]int
}

func NewClient(userAgent string, windowSize [2]int) *Client {
	cookieJar := NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeout(30),
		tls_client.WithClientProfile(tls_client.Chrome_107),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithCookieJar(cookieJar),
		tls_client.WithRandomTLSExtensionOrder(),
	}
	tlsClient, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	return &Client{
		tlsClient:  tlsClient,
		UserAgent:  userAgent,
		WindowSize: windowSize,
	}
}

func (c *Client) NewRequest() *Request {
	cReq := NewRequest(c.tlsClient)
	cReq.SetHeader("Accept", "*/*")
	cReq.SetHeader("Accept-Encoding", "gzip,deflate,br")
	cReq.SetHeader("Cache-Control", "no-store,no-cache")
	cReq.SetHeader("Pragma", "no-cache")
	cReq.SetHeader("User-Agent", c.UserAgent)
	cReq.SetHeaderOrder([]string{"Accept", "Accept-Encoding", "Cache-Control", "Origin", "Pragma", "Referer", "User-Agent"})
	return cReq
}

func (c *Client) SetProxy(proxyUrl string) error {
	c.RawProxy = proxyUrl
	c.ProxyUrl = nil
	var err error
	if proxyUrl != "" {
		if !isUrl(proxyUrl) {
			proxyUrl = "http://" + proxyUrl
		}
		c.ProxyUrl, err = url.Parse(proxyUrl)
		if err != nil {
			return err
		}

	}
	return c.tlsClient.SetProxy(proxyUrl)
}

func (c *Client) SetAutoRedirect(b bool) {
	c.tlsClient.SetFollowRedirect(b)
}

func (c *Client) SetCookies(domain string, cookies []*http.Cookie) {
	u := &url.URL{Host: domain, Scheme: "https", Path: "/"}
	c.tlsClient.SetCookies(u, cookies)
}

func (c *Client) GetRequestInfo() (bool, string) {
	_, b, e := c.NewRequest().SetUrl("https://tls.peet.ws/api/all").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "socks5://")
}
