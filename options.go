package bilibili_go

import "net/http"

type options struct {
	// UserAgent 自定义用户头
	UserAgent string

	// HttpClient 自定义http客户端
	HttpClient *http.Client

	// CookieFilePath cookie缓存文件路径, 如果配置了则会缓存cookie，否则不缓存，默认为空
	CookieFilePath string

	// Debug 是否开启调试模式，如果开启则会将http的请求信息输出到stdout
	Debug bool
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

type debug bool

func (d debug) apply(opt *options) {
	opt.Debug = bool(d)
}

func WithDebug(d bool) Option {
	return debug(d)
}

/* ========================================================== */

var defaultOptions = options{
	UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	HttpClient:     http.DefaultClient,
	CookieFilePath: "",
}

func applyOptions(opts ...Option) *options {
	opt := &defaultOptions
	for _, o := range opts {
		o.apply(opt)
	}

	return opt
}
