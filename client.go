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

type profileList struct {
	term     string
	defaults tls_client.ClientProfile
	profiles map[string]tls_client.ClientProfile
}

var clientProfiles = []profileList{
	{
		term:     "Chrome",
		defaults: tls_client.DefaultClientProfile,
		profiles: map[string]tls_client.ClientProfile{
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
			"Chrome/113": tls_client.Chrome_112,
			"Chrome/114": tls_client.Chrome_112,
			"Chrome/115": tls_client.Chrome_112,
			"Chrome/116": tls_client.Chrome_112,
			"Chrome/117": tls_client.Chrome_117,
			"Chrome/118": tls_client.Chrome_117,
			"Chrome/119": tls_client.Chrome_117,
			"Chrome/120": tls_client.Chrome_120,
		},
	},
	{
		term:     "Firefox",
		defaults: tls_client.Firefox_123,
		profiles: map[string]tls_client.ClientProfile{
			"Firefox/102": tls_client.Firefox_102,
			"Firefox/104": tls_client.Firefox_104,
			"Firefox/105": tls_client.Firefox_105,
			"Firefox/106": tls_client.Firefox_106,
			"Firefox/108": tls_client.Firefox_108,
			"Firefox/109": tls_client.Firefox_108,
			"Firefox/110": tls_client.Firefox_110,
			"Firefox/111": tls_client.Firefox_110,
			"Firefox/112": tls_client.Firefox_110,
			"Firefox/113": tls_client.Firefox_110,
			"Firefox/114": tls_client.Firefox_110,
			"Firefox/115": tls_client.Firefox_110,
			"Firefox/116": tls_client.Firefox_110,
			"Firefox/117": tls_client.Firefox_117,
			"Firefox/118": tls_client.Firefox_117,
			"Firefox/119": tls_client.Firefox_117,
			"Firefox/120": tls_client.Firefox_120,
			"Firefox/121": tls_client.Firefox_120,
			"Firefox/122": tls_client.Firefox_120,
			"Firefox/123": tls_client.Firefox_123,
		},
	},
	{
		term:     "Version",
		defaults: tls_client.Safari_16_0,
		profiles: map[string]tls_client.ClientProfile{
			"Version/15":     tls_client.Safari_15_6_1,
			"iPhone OS 15_5": tls_client.Safari_IOS_15_5,
			"iPhone OS 15_6": tls_client.Safari_IOS_15_6,
			"iPhone OS 16_":  tls_client.Safari_IOS_16_0,
			"iPhone OS 17_":  tls_client.Safari_IOS_17_0,
			"iPad":           tls_client.Safari_Ipad_15_6,
		},
	},
}

var defaultHeaderOrder = []string{
	"content-length",
	"pragma",
	"cache-control",
	"sec-ch-ua",
	"sec-ch-ua-mobile",
	"sec-ch-ua-platform",
	"upgrade-insecure-requests",
	"user-agent",
	"content-type",
	"x-requested-with",
	"accept",
	"origin",
	"host",
	"sec-fetch-site",
	"sec-fetch-mode",
	"sec-fetch-user",
	"sec-fetch-dest",
	"referer",
	"accept-encoding",
	"accept-language",
	"cookie",
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
	_, b, e := c.NewRequest().Get("https://httpbin.org/get").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func (c *Client) GetFingerPrint() (bool, string) {
	_, b, e := c.NewRequest().Get("https://tls.peet.ws/api/all").Send().End()
	if e != nil {
		return false, e.Error()
	}

	return true, b
}

func (c *Client) GetIPLocation() (bool, string) {
	_, b, e := c.NewRequest().Get("http://ip-api.com/json").Send().End()
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
	for _, clientProfile := range clientProfiles {
		if strings.Contains(userAgent, clientProfile.term) {
			profile = clientProfile.defaults
			for k, v := range clientProfile.profiles {
				if strings.Contains(userAgent, k) {
					profile = v
					break
				}
			}
		}
	}
	return profile
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "socks5://")
}
