package bilibili_go

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Client struct {
	httpClient *http.Client
	cookies    []*http.Cookie
	ua         string
}

// TODO Update
func NewClient() *Client {
	return &Client{
		httpClient: http.DefaultClient,
		ua:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36",
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

/* ===================== helper ===================== */

func (c *Client) do(req *http.Request, beforeDo func(request *http.Request), afterDo func(response *http.Response)) (*BaseResponse, error) {
	req.Header.Set("User-Agent", c.ua)

	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	if beforeDo != nil {
		beforeDo(req)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if afterDo != nil {
		afterDo(resp)
	}

	return NewBaseResponse(resp.Body)
}

func (c *Client) get(uri string, param url.Values, beforeDo func(request *http.Request), afterDo func(response *http.Response)) (*BaseResponse, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	if param != nil {
		req.URL.RawQuery = param.Encode()
	}

	return c.do(req, beforeDo, afterDo)
}

func (c *Client) post(uri string, body io.Reader, beforeDo func(request *http.Request), afterDo func(response *http.Response)) (*BaseResponse, error) {
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	return c.do(req, beforeDo, afterDo)
}
