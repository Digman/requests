package requests

import (
	tls_client "github.com/Digman/tls-client"
	http "github.com/bogdanfinn/fhttp"
	"net/url"
	"strings"
)

type Client struct {
	tlsClient tls_client.HttpClient

	HeaderOrder []string
	ProxyUrl    *url.URL
	RawProxy    string
	UserAgent   string
	WindowSize  [2]int
}

var clientProfiles = map[string]map[string]tls_client.ClientProfile{
	"Chrome": {
		"default":    tls_client.DefaultClientProfile,
		"Chrome/103": tls_client.Chrome_103,
		"Chrome/104": tls_client.Chrome_104,
		"Chrome/105": tls_client.Chrome_105,
		"Chrome/106": tls_client.Chrome_106,
		"Chrome/107": tls_client.Chrome_107,
	},
	"iPhone OS": {
		"default":        tls_client.Safari_IOS_15_5,
		"iPhone OS 15_5": tls_client.Safari_IOS_15_5,
		"iPhone OS 15_6": tls_client.Safari_IOS_15_6,
		"iPhone OS 16_0": tls_client.Safari_IOS_16_0,
	},
	"iPad": {
		"default": tls_client.Safari_Ipad_15_6,
	},
}

var defaultHeaderOrder = []string{
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Cache-Control",
	"Content-Length",
	"Content-Type",
	"Cookie",
	"Origin",
	"Pragma",
	"Referer",
	"User-Agent",
	"x-requested-with",
}

func NewClient(userAgent string) *Client {
	clientProfile := getClientProfile(userAgent)
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeout(30),
		tls_client.WithClientProfile(clientProfile),
		tls_client.WithNewCookieJar(),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithRandomTLSExtensionOrder(),
	}
	tlsClient, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	return &Client{
		tlsClient:   tlsClient,
		HeaderOrder: defaultHeaderOrder,
		UserAgent:   userAgent,
		WindowSize:  [2]int{1280, 1080},
	}
}

func (c *Client) NewRequest() *Request {
	cReq := NewRequest(c.tlsClient)
	cReq.SetHeader("Accept", "*/*")
	cReq.SetHeader("Accept-Language", "en-US,en;q=0.9")
	cReq.SetHeader("Accept-Encoding", "gzip,deflate,br")
	cReq.SetHeader("Cache-Control", "max-age=0")
	cReq.SetHeader("Pragma", "no-cache")
	cReq.SetHeader("User-Agent", c.UserAgent)
	cReq.SetHeaderOrder(c.HeaderOrder)
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
	_, b, e := c.NewRequest().SetUrl("https://client.tlsfingerprint.io:8443/").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func getClientProfile(userAgent string) tls_client.ClientProfile {
	profile := tls_client.DefaultClientProfile
	for ua, kv := range clientProfiles {
		if strings.Contains(userAgent, ua) {
			profile = kv["default"]
			for s, clientProfile := range kv {
				if s != "default" && strings.Contains(userAgent, s) {
					profile = clientProfile
					break
				}
			}
			break
		}
	}

	return profile
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "socks5://")
}
