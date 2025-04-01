package requests

import (
	tls_client "github.com/Digman/tls-client"
	"github.com/Digman/tls-client/profiles"
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

type CertPinning struct {
	// CertificatePins host => certificatePins
	CertificatePins map[string][]string
	// BadPinHandler func to handle error
	BadPinHandler tls_client.BadPinHandlerFunc
}

type profileList struct {
	term     string
	defaults profiles.ClientProfile
	profiles map[string]profiles.ClientProfile
}

var clientProfiles = []profileList{
	{
		term:     "Chrome",
		defaults: profiles.DefaultClientProfile,
		profiles: map[string]profiles.ClientProfile{
			"Chrome/103": profiles.Chrome_103,
			"Chrome/104": profiles.Chrome_104,
			"Chrome/105": profiles.Chrome_105,
			"Chrome/106": profiles.Chrome_106,
			"Chrome/107": profiles.Chrome_107,
			"Chrome/108": profiles.Chrome_108,
			"Chrome/109": profiles.Chrome_109,
			"Chrome/110": profiles.Chrome_110,
			"Chrome/111": profiles.Chrome_111,
			"Chrome/112": profiles.Chrome_112,
			"Chrome/113": profiles.Chrome_112_PSK,
			"Chrome/114": profiles.Chrome_114_PSK,
			"Chrome/115": profiles.Chrome_114_PSK,
			"Chrome/116": profiles.Chrome_116_PSK,
			"Chrome/117": profiles.Chrome_117_PSK,
			"Chrome/118": profiles.Chrome_117_PSK,
			"Chrome/119": profiles.Chrome_120_PSK,
			"Chrome/120": profiles.Chrome_120_PSK,
			"Chrome/121": profiles.Chrome_120_PSK,
			"Chrome/122": profiles.Chrome_120_PSK,
			"Chrome/123": profiles.Chrome_120_PSK,
			"Chrome/124": profiles.Chrome_124_PSK,
			"Chrome/125": profiles.Chrome_124_PSK,
			"Chrome/126": profiles.Chrome_124_PSK,
			"Chrome/127": profiles.Chrome_124_PSK,
			"Chrome/128": profiles.Chrome_124_PSK,
			"Chrome/129": profiles.Chrome_124_PSK,
			"Chrome/130": profiles.Chrome_124_PSK,
			"Chrome/131": profiles.Chrome_131_PSK,
			"Chrome/133": profiles.Chrome_133_PSK,
		},
	},
	{
		term:     "Firefox",
		defaults: profiles.Firefox_135,
		profiles: map[string]profiles.ClientProfile{
			"Firefox/102": profiles.Firefox_102,
			"Firefox/104": profiles.Firefox_104,
			"Firefox/105": profiles.Firefox_105,
			"Firefox/106": profiles.Firefox_106,
			"Firefox/108": profiles.Firefox_108,
			"Firefox/109": profiles.Firefox_108,
			"Firefox/110": profiles.Firefox_110,
			"Firefox/111": profiles.Firefox_110,
			"Firefox/112": profiles.Firefox_110,
			"Firefox/113": profiles.Firefox_110,
			"Firefox/114": profiles.Firefox_117,
			"Firefox/115": profiles.Firefox_117,
			"Firefox/116": profiles.Firefox_117,
			"Firefox/117": profiles.Firefox_117,
			"Firefox/118": profiles.Firefox_117,
			"Firefox/119": profiles.Firefox_117,
			"Firefox/120": profiles.Firefox_123,
			"Firefox/121": profiles.Firefox_123,
			"Firefox/122": profiles.Firefox_123,
			"Firefox/123": profiles.Firefox_123,
			"Firefox/132": profiles.Firefox_132,
			"Firefox/133": profiles.Firefox_133,
			"Firefox/135": profiles.Firefox_135,
		},
	},
	{
		term:     "Version",
		defaults: profiles.Safari_15_6_1,
		profiles: map[string]profiles.ClientProfile{
			"Version/15":     profiles.Safari_15_6_1,
			"Version/16":     profiles.Safari_16_0,
			"iPhone OS 15_5": profiles.Safari_IOS_15_5,
			"iPhone OS 15_6": profiles.Safari_IOS_15_6,
			"iPhone OS 16_":  profiles.Safari_IOS_16_0,
			"iPhone OS 17_":  profiles.Safari_IOS_17_0,
			"iPhone OS 18_":  profiles.Safari_IOS_18_0,
			"iPad":           profiles.Safari_Ipad_15_6,
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

// defaultUserAgent default useragent
var defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36"

// defaultWindowSize default window size
var defaultWindowSize = [2]int{1440, 900}

// defaultTimeout default request timeout in milliseconds
var defaultTimeout = 30000

func NewClient(userAgent string, cp ...*CertPinning) *Client {
	return newClient(userAgent, defaultWindowSize, defaultTimeout, cp...)
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

func (c *Client) SetKeepAlive(b bool) {
	if c.ExtraHeaders == nil {
		c.ExtraHeaders = make(map[string]string)
	}
	if b {
		c.ExtraHeaders["Connection"] = "keep-alive"
	} else {
		c.ExtraHeaders["Connection"] = "close"
	}
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

func (c *Client) SetUrlCookies(cookieUrl string, cookies []*http.Cookie) {
	if u, err := url.Parse(cookieUrl); err == nil {
		c.tlsClient.SetCookies(u, cookies)
	}
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

func newClient(userAgent string, windowSize [2]int, timeout int, cp ...*CertPinning) *Client {
	clientProfile := getClientProfile(userAgent)
	cookieJar := tls_client.NewCookieJar()
	connectHeader := http.Header{"User-Agent": {userAgent}}
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutMilliseconds(timeout),
		tls_client.WithClientProfile(clientProfile),
		tls_client.WithCookieJar(cookieJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithConnectHeaders(connectHeader),
	}
	if len(cp) > 0 && cp[0] != nil {
		options = append(options, tls_client.WithCertificatePinning(
			cp[0].CertificatePins, cp[0].BadPinHandler),
		)
	}
	tlsClient, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	return &Client{
		tlsClient:   tlsClient,
		HeaderOrder: defaultHeaderOrder,
		UserAgent:   userAgent,
		WindowSize:  windowSize,
	}
}

func getClientProfile(userAgent string) profiles.ClientProfile {
	for _, clientProfile := range clientProfiles {
		if strings.Contains(userAgent, clientProfile.term) {
			for k, v := range clientProfile.profiles {
				if strings.Contains(userAgent, k) {
					return v
				}
			}
			return clientProfile.defaults
		}
	}
	return profiles.DefaultClientProfile
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "socks5://")
}
