package bilibili_go

import (
	"fmt"
	"github.com/kainhuck/bilibili-go/internal/net"
	"github.com/kainhuck/bilibili-go/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	httpClient       *net.HttpClient
	authInfo         *authInfo
	cookieCache      map[string]string
	authFilePath     string
	wbiKey           string
	wbiKeyLastUpdate time.Time
	debug            bool
}

func NewClient(opts ...Option) *Client {
	opt := applyOptions(opts...)

	auth, _ := loadAuthInfoFromFile(opt.AuthFilePath)

	client := net.NewHttpClient(opt.HttpClient).SetUserAgent(opt.UserAgent)

	return &Client{
		httpClient:   client,
		authInfo:     auth,
		cookieCache:  make(map[string]string),
		authFilePath: opt.AuthFilePath,
		debug:        opt.Debug,
	}
}

func (c *Client) setAuthInfo(auth *authInfo) {
	c.authInfo = auth
	if c.authInfo == nil {
		return
	}
	for _, cookie := range c.authInfo.Cookies {
		c.cookieCache[cookie.Name] = cookie.Value
	}
}

func (c *Client) loadAuthInfoFromFile() bool {
	if utils.FileExists(c.authFilePath) {
		auth, err := loadAuthInfoFromFile(c.authFilePath)
		if err != nil {
			logrus.Errorf("load auth info from file: %v failed: %v", c.authFilePath, err)
			return false
		}

		if auth != nil {
			c.setAuthInfo(auth)
			logrus.Infof("load auth info from: %v", c.authFilePath)
			return true
		}
	}

	return false
}

func (c *Client) getWbiKey() string {
	resp, err := c.GetNavigation()
	if err != nil {
		return ""
	}

	imgKey := strings.Split(strings.Split(resp.WBIImg.ImgURL, "/")[len(strings.Split(resp.WBIImg.ImgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(resp.WBIImg.SubURL, "/")[len(strings.Split(resp.WBIImg.SubURL, "/"))-1], ".")[0]
	return getMixinKey(imgKey + subKey)
}

func (c *Client) updateWbiKeyCache() {
	if time.Since(c.wbiKeyLastUpdate).Minutes() < 10 {
		return
	}

	c.wbiKey = c.getWbiKey()
	c.wbiKeyLastUpdate = time.Now()
}

func (c *Client) getWbiKeyCached() string {
	c.updateWbiKeyCache()

	return c.wbiKey
}

// LoginWithQrCode 登陆这一步必须成功，否则后续接口无法访问
func (c *Client) LoginWithQrCode() {
	if c.loadAuthInfoFromFile() {
		return
	}

	defer func() { _ = saveAuthInfoToFile(c.authFilePath, c.authInfo) }()

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
			c.authInfo.RefreshToken = resp.RefreshToken
			logrus.Infof("login success!!!")
			return
		case 86038:
			logrus.Errorf("qrcode expired")
			os.Exit(-1)
		}
		time.Sleep(1 * time.Second)
	}
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

	if c.debug {
		client = client.Debug()
	}

	if auth && c.authInfo != nil {
		client = client.SetCookies(c.authInfo.Cookies)
	}

	return client
}
