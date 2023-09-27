package requests

import (
	tls_client "github.com/Digman/tls-client"
	http "github.com/bogdanfinn/fhttp"
	"net/url"
	"strings"
)

type Client struct {
	tlsClient tls_client.HttpClient

	ExtraHeaders map[string]string
	HeaderOrder  []string
	ProxyUrl     *url.URL
	RawProxy     string
	UserAgent    string
	WindowSize   [2]int
}

var clientProfiles = map[string]map[string]tls_client.ClientProfile{
	"Chrome": {
		"default":    tls_client.DefaultClientProfile,
		"Chrome/100": tls_client.Chrome_100,
		"Chrome/101": tls_client.Chrome_100,
		"Chrome/102": tls_client.Chrome_102,
		"Chrome/103": tls_client.Chrome_103,
		"Chrome/104": tls_client.Chrome_104,
		"Chrome/105": tls_client.Chrome_105,
		"Chrome/106": tls_client.Chrome_106,
		"Chrome/107": tls_client.Chrome_107,
		"Chrome/108": tls_client.Chrome_108,
		"Chrome/109": tls_client.Chrome_109,
		"Chrome/110": tls_client.Chrome_110,
		"Chrome/111": tls_client.Chrome_111,
		"Chrome/112": tls_client.Chrome_112,
	},
	"iPhone": {
		"default": tls_client.Safari_IOS_16_0,
		"OS 15_5": tls_client.Safari_IOS_15_5,
		"OS 15_6": tls_client.Safari_IOS_15_6,
		"OS 15_7": tls_client.Safari_IOS_15_6,
		"OS 16_0": tls_client.Safari_IOS_16_0,
		"OS 16_1": tls_client.Safari_IOS_16_0,
		"OS 16_2": tls_client.Safari_IOS_16_0,
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
	"Cookie",
	"Pragma",
	"Sec-Ch-Ua",
	"Sec-Ch-Ua-Mobile",
	"Sec-Ch-Ua-Platform",
	"Sec-Fetch-Dest",
	"Sec-Fetch-Mode",
	"Sec-Fetch-Site",
	"Sec-Fetch-User",
	"Connection",
	"Content-Length",
	"Content-Type",
	"Origin",
	"Referer",
	"User-Agent",
	"X-Requested-With",
}

var defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36"

var defaultWindowSize = [2]int{1440, 900}

var defaultTimeout = 30000

func NewClient(userAgent string) *Client {
	return newClient(userAgent, defaultWindowSize, defaultTimeout)
}

func TimeoutClient(timeout int) *Client {
	return newClient(defaultUserAgent, defaultWindowSize, timeout)
}

func DefaultClient() *Client {
	return NewClient(defaultUserAgent)
}

func (c *Client) NewRequest() *Request {
	cReq := NewRequest(c.tlsClient)
	cReq.SetHeader("Accept", "*/*")
	cReq.SetHeader("Accept-Language", "en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7")
	cReq.SetHeader("Accept-Encoding", "gzip,deflate,br")
	cReq.SetHeader("Cache-Control", "no-cache")
	// cReq.SetHeader("Connection", "keep-alive")
	cReq.SetHeader("Pragma", "no-cache")
	cReq.SetHeader("User-Agent", c.UserAgent)
	if len(c.ExtraHeaders) > 0 {
		for k, v := range c.ExtraHeaders {
			cReq.SetHeader(k, v)
		}
	}
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

func (c *Client) GetCookies(domain string) []*http.Cookie {
	u := &url.URL{Host: domain, Scheme: "https", Path: "/"}
	return c.tlsClient.GetCookies(u)
}

func (c *Client) GetUrlCookies(cookieUrl string) []*http.Cookie {
	if u, err := url.Parse(cookieUrl); err == nil {
		return c.tlsClient.GetCookies(u)
	}
	return nil
}

func (c *Client) GetRequestInfo() (bool, string) {
	_, b, e := c.NewRequest().SetUrl("http://httpbin.org/get").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func (c *Client) GetFingerPrint() (bool, string) {
	_, b, e := c.NewRequest().SetUrl("https://tls.peet.ws/api/all").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func newClient(userAgent string, windowSize [2]int, timeout int) *Client {
	clientProfile := getClientProfile(userAgent)
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeout(2592000),        // Transport timeout(Second)
		tls_client.WithRequestTimeout(timeout), // Request timeout(Millisecond)
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
		WindowSize:  windowSize,
	}
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
