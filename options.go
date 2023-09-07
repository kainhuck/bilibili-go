package bilibili_go

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
)

type options struct {
	// UserAgent 自定义用户头
	UserAgent string

	// HttpClient 自定义http客户端
	HttpClient *http.Client

	// AuthStorage 认证信息存储
	AuthStorage AuthStorage

	// Debug 是否开启调试模式，如果开启则会将http的请求信息输出到output，如果output为nil则视为os.Stdout
	Debug *debugInfo

	// Logger 自定义日志
	Logger Logger

	// ShowQRCodeFunc 自定义输出二维码的方法，默认在stdout输出，
	// 可以通过自定义该方法可以实现其他输出，比如将图片发送到消息通知群
	ShowQRCodeFunc func(code *qrcode.QRCode) error
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

type authStorage struct {
	storage AuthStorage
}

func (a authStorage) apply(opt *options) {
	opt.AuthStorage = a.storage
}

func WithAuthStorage(storage AuthStorage) Option {
	return authStorage{storage: storage}
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

type showQRCodeFunc func(code *qrcode.QRCode) error

func (s showQRCodeFunc) apply(opt *options) {
	opt.ShowQRCodeFunc = s
}

func WithShowQRCodeFunc(f func(code *qrcode.QRCode) error) Option {
	return showQRCodeFunc(f)
}

/* ========================================================== */

var defaultOptions = options{
	UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
	HttpClient: http.DefaultClient,
	Debug:      &debugInfo{},
	Logger:     logrus.StandardLogger(),
	ShowQRCodeFunc: func(code *qrcode.QRCode) error {
		_, err := fmt.Fprint(os.Stdout, code.ToSmallString(true))

		return err
	},
}

func applyOptions(opts ...Option) *options {
	opt := &defaultOptions
	for _, o := range opts {
		o.apply(opt)
	}

	return opt
}
