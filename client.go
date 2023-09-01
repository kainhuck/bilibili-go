package bilibili_go

import (
	"fmt"
	"github.com/kainhuck/bilibili-go/internal/net"
	"github.com/kainhuck/bilibili-go/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	httpClient     *net.HttpClient
	cookies        []*http.Cookie
	cookieCache    map[string]string
	cookieFilePath string
}

func NewClient(opts ...Option) *Client {
	opt := applyOptions(opts...)

	cookies, _ := utils.LoadCookiesFromFile(opt.CookieFilePath)

	client := net.NewHttpClient(opt.HttpClient).SetUserAgent(opt.UserAgent)

	return &Client{
		httpClient:     client,
		cookies:        cookies,
		cookieCache:    make(map[string]string),
		cookieFilePath: opt.CookieFilePath,
	}
}

func (c *Client) setCookies(cookies []*http.Cookie) {
	c.cookies = cookies
	for _, cookie := range cookies {
		c.cookieCache[cookie.Name] = cookie.Value
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
			c.setCookies(cookies)
			logrus.Infof("load cookie from: %v", c.cookieFilePath)
			return
		}
	}

	c.LoginWithQrCode()

	_ = utils.SaveCookiesToFile(c.cookieFilePath, c.cookies)
}

// UploadVideo 视频上传 videoPath 视频路径
func (c *Client) UploadVideo(videoPath string) (*Video, error) {
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return nil, err
	}

	// 调接口上传
	// 1. 预上传
	preResp, err := c.preUpload(fileInfo.Name(), fileInfo.Size())
	if err != nil {
		return nil, err
	}
	if preResp.OK != 1 {
		return nil, fmt.Errorf("[preUpload] upload failed code: %v", preResp.OK)
	}

	// 2. 获取 upload_id
	uploadIDResp, err := c.getUploadID(preResp.Uri(), preResp.Auth, preResp.BizID, fileInfo.Size())
	if err != nil {
		return nil, err
	}
	if uploadIDResp.OK != 1 {
		return nil, fmt.Errorf("[getUploadID] upload failed code: %v", uploadIDResp.OK)
	}

	// 3. 分片上传
	total, err := os.ReadFile(videoPath)
	if err != nil {
		return nil, err
	}

	// 分区
	parts := utils.SplitBytes(total, 10*utils.MB)

	errChan := make(chan error, len(parts))
	wg := sync.WaitGroup{}

	wg.Add(len(parts))
	go func() {
		wg.Wait()
		close(errChan)
	}()
	for i, part := range parts {
		go func(part []byte, number int) {
			defer wg.Done()
			errChan <- c.uploadFileClip(preResp.Uri(), preResp.Auth, uploadIDResp.UploadID, number, len(parts), len(part), (number-1)*10*utils.MB, (number-1)*10*utils.MB+len(part), fileInfo.Size(), part)
		}(part, i+1)
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// 视频上传完成
	checkResp, err := c.uploadCheck(preResp.Uri(), preResp.Auth, fileInfo.Name(), uploadIDResp.UploadID, preResp.BizID)
	if err != nil {
		return nil, err
	}

	if checkResp.OK != 1 {
		return nil, fmt.Errorf("[uploadCheck] upload failed code: %v", checkResp.OK)
	}

	return &Video{
		Filename: preResp.Filename(),
		Title:    strings.Split(fileInfo.Name(), ".")[0],
		Desc:     "",
		CID:      preResp.BizID,
	}, nil
}

/* ===================== helper ===================== */

func (c *Client) getHttpClient(auth bool) *net.HttpClient {
	client := c.httpClient.Clone()

	if auth {
		client.SetCookies(c.cookies)
	}

	return client
}
