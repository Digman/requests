package requests

import (
	"maps"
	"net/url"
	"strings"
	"time"

	tls_client "github.com/Digman/tls-client"
	"github.com/Digman/tls-client/profiles"
	http "github.com/bogdanfinn/fhttp"
)

type Client struct {
	tlsClient tls_client.HttpClient

	ExtraHeaders map[string]string
	HeaderOrder  []string
	ProxyUrl     *url.URL
	RawProxy     string
	UserAgent    string
	WindowSize   [2]int

	// 預組裝請求頭模板，避免每次 NewRequest 重複 Set
	headerTemplate map[string]string
}

type CertPinning struct {
	// CertificatePins host => certificatePins
	CertificatePins map[string][]string
	// BadPinHandler func to handle error
	BadPinHandler tls_client.BadPinHandlerFunc
}

type profileList struct {
	// terms[0] 同時是 profile map key 的前綴；後續元素為等價別名，UA 含任一即命中，比對時替換 key 前綴
	terms    []string
	defaults profiles.ClientProfile
	profiles map[string]profiles.ClientProfile
}

var clientProfiles = []profileList{
	{
		// CriOS = iOS 上的 Chrome，TLS 指紋與桌面 Chrome 同步
		terms:    []string{"Chrome", "CriOS"},
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
			"Chrome/130": profiles.Chrome_130_PSK,
			"Chrome/131": profiles.Chrome_131_PSK,
			"Chrome/132": profiles.Chrome_132_PSK,
			"Chrome/133": profiles.Chrome_133_PSK,
			"Chrome/134": profiles.Chrome_134_PSK,
			"Chrome/135": profiles.Chrome_135_PSK,
			"Chrome/136": profiles.Chrome_136_PSK,
			"Chrome/137": profiles.Chrome_137_PSK,
			"Chrome/138": profiles.Chrome_138_PSK,
			"Chrome/139": profiles.Chrome_139_PSK,
			"Chrome/140": profiles.Chrome_140_PSK,
			"Chrome/141": profiles.Chrome_141_PSK,
			"Chrome/142": profiles.Chrome_142_PSK,
			"Chrome/143": profiles.Chrome_143_PSK,
			"Chrome/144": profiles.Chrome_144_PSK,
			"Chrome/145": profiles.Chrome_145_PSK,
			"Chrome/146": profiles.Chrome_146_PSK,
		},
	},
	{
		terms:    []string{"Firefox"},
		defaults: profiles.Firefox_147,
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
			"Firefox/136": profiles.Firefox_135,
			"Firefox/137": profiles.Firefox_135,
			"Firefox/138": profiles.Firefox_135,
			"Firefox/139": profiles.Firefox_135,
			"Firefox/140": profiles.Firefox_135,
			"Firefox/141": profiles.Firefox_135,
			"Firefox/142": profiles.Firefox_135,
			"Firefox/143": profiles.Firefox_135,
			"Firefox/144": profiles.Firefox_135,
			"Firefox/145": profiles.Firefox_135,
			"Firefox/146": profiles.Firefox_146_PSK,
			"Firefox/147": profiles.Firefox_147,
		},
	},
	{
		terms:    []string{"Version"},
		defaults: profiles.Safari_26,
		profiles: map[string]profiles.ClientProfile{
			"Version/15":     profiles.Safari_15_6_1,
			"Version/16":     profiles.Safari_16_0,
			"Version/17":     profiles.Safari_IOS_17_0,
			"Version/18":     profiles.Safari_IOS_18_5,
			"Version/26":     profiles.Safari_26,
			"iPhone OS 15_5": profiles.Safari_IOS_15_5,
			"iPhone OS 15_6": profiles.Safari_IOS_15_6,
			"iPhone OS 16_":  profiles.Safari_IOS_16_0,
			"iPhone OS 17_":  profiles.Safari_IOS_17_0,
			"iPhone OS 18_0": profiles.Safari_IOS_18_0,
			"iPhone OS 18_":  profiles.Safari_IOS_26_0,
			"CPU OS 17_":     profiles.Safari_IOS_17_0,
			"CPU OS 18_0":    profiles.Safari_IOS_18_0,
			"CPU OS 18_":     profiles.Safari_IOS_26_0,
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
var defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"

// defaultWindowSize default window size
var defaultWindowSize = [2]int{1440, 900}

// defaultTimeout default request timeout in milliseconds
var defaultTimeout = 30000

// PoolConfig 連線池容量設定，比直接暴露 tls_client.TransportOptions 更語意化
type PoolConfig struct {
	MaxIdleConns        int           // 全局 idle 連線上限（跨 host）
	MaxIdleConnsPerHost int           // 單 host idle 連線上限
	MaxConnsPerHost     int           // 單 host 總連線上限
	IdleConnTimeout     time.Duration // idle 連線存活時間
}

// 預設連線池組合，依使用場景選用
var (
	// PoolDefault 通用場景（4-8 worker 共享 Client）
	PoolDefault = PoolConfig{MaxIdleConns: 100, MaxIdleConnsPerHost: 32, MaxConnsPerHost: 64, IdleConnTimeout: 90 * time.Second}

	// PoolSmall 每 task 獨立 Client 的高並發場景（如 500+ 抢购 task）
	// 單 Client 池小，避免整體 fd/內存暴增
	PoolSmall = PoolConfig{MaxIdleConns: 16, MaxIdleConnsPerHost: 4, MaxConnsPerHost: 16, IdleConnTimeout: 60 * time.Second}

	// PoolLarge 共享 Client 且超高並發
	PoolLarge = PoolConfig{MaxIdleConns: 500, MaxIdleConnsPerHost: 128, MaxConnsPerHost: 256, IdleConnTimeout: 120 * time.Second}
)

func (p PoolConfig) toTransport() *tls_client.TransportOptions {
	idle := p.IdleConnTimeout
	return &tls_client.TransportOptions{
		MaxIdleConns:        p.MaxIdleConns,
		MaxIdleConnsPerHost: p.MaxIdleConnsPerHost,
		MaxConnsPerHost:     p.MaxConnsPerHost,
		IdleConnTimeout:     &idle,
	}
}

func NewClient(userAgent string, cp ...*CertPinning) *Client {
	return newClient(userAgent, defaultWindowSize, defaultTimeout, PoolDefault, cp...)
}

// NewClientWithPool 自定義連線池容量建立 Client，適用於高並發或共享場景
//
//	小池（500+ task 各持一個 Client）: requests.NewClientWithPool(ua, requests.PoolSmall)
//	自訂: requests.NewClientWithPool(ua, requests.PoolConfig{MaxIdleConnsPerHost: 8, ...})
func NewClientWithPool(userAgent string, pool PoolConfig, cp ...*CertPinning) *Client {
	return newClient(userAgent, defaultWindowSize, defaultTimeout, pool, cp...)
}

func TimeoutClient(timeout int) *Client {
	return newClient(defaultUserAgent, defaultWindowSize, timeout, PoolDefault)
}

func DefaultClient() *Client {
	return NewClient(defaultUserAgent)
}

func (c *Client) NewRequest() *Request {
	cReq := NewRequest(c.tlsClient)
	cReq.header = maps.Clone(c.headerTemplate)
	if len(c.ExtraHeaders) > 0 {
		maps.Copy(cReq.header, c.ExtraHeaders)
	}
	cReq.SetHeaderOrder(c.HeaderOrder)
	return cReq
}

// buildHeaderTemplate 在 Client 構造時預組裝通用請求頭，避免每次 NewRequest 重複賦值
func (c *Client) buildHeaderTemplate() {
	c.headerTemplate = map[string]string{
		"Accept":          "*/*",
		"Accept-Language": "en-US,en;q=0.9,zh-TW;q=0.8,zh;q=0.7",
		"Accept-Encoding": "gzip, deflate, br",
		"Cache-Control":   "no-cache",
		"Pragma":          "no-cache",
		"User-Agent":      c.UserAgent,
	}
}

func (c *Client) SetKeepAlive(b bool) {
	if c.ExtraHeaders == nil {
		c.ExtraHeaders = make(map[string]string)
	}
	if b {
		// 恢復默認保活：刪除顯式 close 標記（Chrome H2 不發 Connection 頭，H1 默認 keep-alive）
		delete(c.ExtraHeaders, "Connection")
	} else {
		// fhttp 的 isConnectionCloseRequest 會識別此 header 並在 stream 結束後關連線
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

func (c *Client) NewCookies() {
	c.tlsClient.SetCookieJar(tls_client.NewCookieJar())
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

func newClient(userAgent string, windowSize [2]int, timeout int, pool PoolConfig, cp ...*CertPinning) *Client {
	clientProfile := getClientProfile(userAgent)
	cookieJar := tls_client.NewCookieJar()
	connectHeader := http.Header{"User-Agent": {userAgent}}
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutMilliseconds(timeout),
		tls_client.WithClientProfile(clientProfile),
		tls_client.WithCookieJar(cookieJar),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithConnectHeaders(connectHeader),
		tls_client.WithTransportOptions(pool.toTransport()),
	}
	// Chrome 使用 extension 隨機排列（與真實 Chrome SSL_set_permute_extensions 一致）
	// Safari/Firefox 不做 extension 隨機排列，順序固定
	if !strings.Contains(userAgent, "Version/") {
		options = append(options, tls_client.WithRandomTLSExtensionOrder())
	}
	if len(cp) > 0 && cp[0] != nil {
		options = append(options, tls_client.WithCertificatePinning(
			cp[0].CertificatePins, cp[0].BadPinHandler),
		)
	}
	tlsClient, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	c := &Client{
		tlsClient:   tlsClient,
		HeaderOrder: defaultHeaderOrder,
		UserAgent:   userAgent,
		WindowSize:  windowSize,
	}
	c.buildHeaderTemplate()
	return c
}

func getClientProfile(userAgent string) profiles.ClientProfile {
	for _, pl := range clientProfiles {
		hit := ""
		for _, t := range pl.terms {
			if strings.Contains(userAgent, t) {
				hit = t
				break
			}
		}
		if hit == "" {
			continue
		}
		canonical := pl.terms[0]
		for k, v := range pl.profiles {
			matchKey := strings.Replace(k, canonical, hit, 1)
			if strings.Contains(userAgent, matchKey) {
				return v
			}
		}
		return pl.defaults
	}
	return profiles.DefaultClientProfile
}

func isUrl(path string) bool {
	return strings.HasPrefix(path, "http://") ||
		strings.HasPrefix(path, "https://") ||
		strings.HasPrefix(path, "socks5://")
}
