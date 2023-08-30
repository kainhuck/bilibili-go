package bilibili_go

import (
	"fmt"
	"github.com/kainhuck/bilibili-go/internal/net"
	"github.com/kainhuck/bilibili-go/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
	"time"
)

type Client struct {
	httpClient     *net.HttpClient
	cookies        []*http.Cookie
	ua             string
	cookieFilePath string
}

func NewClient(opts ...Option) *Client {
	opt := applyOptions(opts...)

	cookies, _ := utils.LoadCookiesFromFile(opt.CookieFilePath)

	return &Client{
		httpClient:     net.NewHttpClient(opt.HttpClient),
		cookies:        cookies,
		ua:             opt.UserAgent,
		cookieFilePath: opt.CookieFilePath,
	}
}

// LoginWithQrCode 登陆这一步必须成功，否则后续接口无法访问
func (c *Client) LoginWithQrCode() {
	generateResp, err := c.qrcodeGenerate()
	if err != nil {
		logrus.Errorf("generate qrcode failed")
		os.Exit(-1)
	}

	qrCode, err := qrcode.New(generateResp.Url, qrcode.Medium)
	if err != nil {
		logrus.Errorf("new qrcode failed")
		os.Exit(-1)
	}

	_, err = fmt.Fprint(os.Stdout, qrCode.ToSmallString(true))
	if err != nil {
		logrus.Errorf("print qrcode failed")
		os.Exit(-1)
	}

	for {
		resp, err := c.qrcodePoll(generateResp.QrcodeKey)
		if err != nil {
			logrus.Errorf("poll qrcode failed")
			os.Exit(-1)
		}

		switch resp.Code {
		case 0:
			logrus.Infof("login success!!!")
			return
		case 86038:
			logrus.Errorf("qrcode expired")
			os.Exit(-1)
		}
		time.Sleep(1 * time.Second)
	}
}

// LoginWithQrCodeWithCache 带有缓存的二维码登陆
func (c *Client) LoginWithQrCodeWithCache() {

	if utils.FileExists(c.cookieFilePath) {
		cookies, err := utils.LoadCookiesFromFile(c.cookieFilePath)
		if err != nil {
			logrus.Errorf("load cookie from file: %v failed: %v", c.cookieFilePath, err)
			os.Exit(-1)
		}

		if len(cookies) != 0 {
			c.cookies = cookies
			logrus.Infof("load cookie from: %v", c.cookieFilePath)
			return
		}
	}

	c.LoginWithQrCode()

	_ = utils.SaveCookiesToFile(c.cookieFilePath, c.cookies)
}

/* ===================== helper ===================== */

func (c *Client) getHttpClient(auth bool) *net.HttpClient {
	client := c.httpClient.Clone()

	if auth {
		client.SetCookies(c.cookies)
	}

	return client
}
