package bilibili_go

import "net/http"

type options struct {
	// UserAgent 自定义用户头
	UserAgent string

	// HttpClient 自定义http客户端
	HttpClient *http.Client

	// CookieFilePath cookie缓存文件路径
	CookieFilePath string
}

type Option interface {
	apply(*options)
}

type userAgent string

func (ua userAgent) apply(opt *options) {
	opt.UserAgent = string(ua)
}

func WithUserAgent(ua string) Option {
	return userAgent(ua)
}

type httpClient struct {
	client *http.Client
}

func (client httpClient) apply(opt *options) {
	opt.HttpClient = client.client
}

func WithHttpClient(client *http.Client) Option {
	return httpClient{client: client}
}

type cookieFilePath string

func (path cookieFilePath) apply(opt *options) {
	opt.CookieFilePath = string(path)
}

func WithCookieFilePath(path string) Option {
	return cookieFilePath(path)
}

/* ========================================================== */

var defaultOptions = options{
	UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	HttpClient:     http.DefaultClient,
	CookieFilePath: "bilibili_cookie.txt",
}

func applyOptions(opts ...Option) *options {
	opt := &defaultOptions
	for _, o := range opts {
		o.apply(opt)
	}

	return opt
}
