package bilibili_go

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type options struct {
	// UserAgent 自定义用户头
	UserAgent string

	// HttpClient 自定义http客户端
	HttpClient *http.Client

	// AuthFilePath cookie缓存文件路径, 如果配置了则会缓存cookie，否则不缓存，默认为空
	AuthFilePath string

	// Debug 是否开启调试模式，如果开启则会将http的请求信息输出到output，如果output为nil则视为os.Stdout
	Debug *debugInfo

	// Logger 自定义日志
	Logger Logger
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

type authFilePath string

func (path authFilePath) apply(opt *options) {
	opt.AuthFilePath = string(path)
}

func WithAuthFilePath(path string) Option {
	return authFilePath(path)
}

type debug struct {
	debugInfo *debugInfo
}

func (d debug) apply(opt *options) {
	opt.Debug = d.debugInfo
}

func WithDebug(d bool, output ...*os.File) Option {
	info := &debugInfo{
		debug: d,
	}

	if len(output) != 0 {
		info.output = output[0]
	}
	return debug{info}
}

type log struct {
	logger Logger
}

func (l log) apply(opt *options) {
	opt.Logger = l.logger
}

func WithLogger(logger Logger) Option {
	return log{logger: logger}
}

/* ========================================================== */

var defaultOptions = options{
	UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	HttpClient:   http.DefaultClient,
	AuthFilePath: "",
	Debug:        &debugInfo{},
	Logger:       logrus.StandardLogger(),
}

func applyOptions(opts ...Option) *options {
	opt := &defaultOptions
	for _, o := range opts {
		o.apply(opt)
	}

	return opt
}
