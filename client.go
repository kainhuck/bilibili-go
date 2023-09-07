package bilibili_go

import (
	"fmt"
	"github.com/kainhuck/bilibili-go/internal/net"
	"github.com/kainhuck/bilibili-go/internal/utils"
	"github.com/skip2/go-qrcode"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type debugInfo struct {
	debug  bool
	output *os.File
}

type Client struct {
	httpClient       *net.HttpClient
	authInfo         *authInfo
	cookieCache      map[string]string
	authFilePath     string
	wbiKey           string
	wbiKeyLastUpdate time.Time
	debug            *debugInfo
	logger           Logger
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
		logger:       opt.Logger,
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
			c.logger.Errorf("load auth info from file: %v failed: %v", c.authFilePath, err)
			return false
		}

		if auth != nil {
			c.setAuthInfo(auth)
			c.logger.Infof("load auth info from: %v", c.authFilePath)
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
		c.logger.Errorf("generate qrcode failed")
		os.Exit(-1)
	}

	qrCode, err := qrcode.New(generateResp.Url, qrcode.Medium)
	if err != nil {
		c.logger.Errorf("new qrcode failed")
		os.Exit(-1)
	}

	_, err = fmt.Fprint(os.Stdout, qrCode.ToSmallString(true))
	if err != nil {
		c.logger.Errorf("print qrcode failed")
		os.Exit(-1)
	}

	for {
		resp, cookies, err := c.qrcodePoll(generateResp.QrcodeKey)
		if err != nil {
			c.logger.Errorf("poll qrcode failed")
			os.Exit(-1)
		}

		switch resp.Code {
		case 0:
			c.setAuthInfo(&authInfo{
				Cookies:      cookies,
				RefreshToken: resp.RefreshToken,
			})
			c.logger.Infof("login success!!!")
			return
		case 86038:
			c.logger.Errorf("qrcode expired")
			os.Exit(-1)
		}
		time.Sleep(1 * time.Second)
	}
}

// UploadVideoFromDisk 从本地磁盘上传视频 videoPath 视频路径
func (c *Client) UploadVideoFromDisk(videoPath string) (*Video, error) {
	fileInfo, err := os.Stat(videoPath)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(videoPath)
	if err != nil {
		return nil, err
	}

	return c.UploadVideo(fileInfo.Name(), content)
}

// UploadVideoFromReader ...
func (c *Client) UploadVideoFromReader(filename string, reader io.Reader) (*Video, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return c.UploadVideo(filename, content)
}

// UploadVideoFromHTTP 从http链接上传文件
func (c *Client) UploadVideoFromHTTP(filename string, url string) (*Video, error) {
	c.logger.Infof("start download file: %v, from: %v", filename, url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.UploadVideoFromReader(filename, resp.Body)
}

// UploadVideo 视频上传，filename 文件名 content 视频内容
func (c *Client) UploadVideo(filename string, content []byte) (*Video, error) {
	filesize := int64(len(content))
	// 调接口上传
	// 1. 预上传
	preResp, err := c.preUpload(filename, filesize)
	if err != nil {
		return nil, err
	}
	if preResp.OK != 1 {
		return nil, fmt.Errorf("[preUpload] upload failed code: %v", preResp.OK)
	}

	// 2. 获取 upload_id
	uploadIDResp, err := c.getUploadID(preResp.Uri(), preResp.Auth, preResp.BizID, filesize)
	if err != nil {
		return nil, err
	}
	if uploadIDResp.OK != 1 {
		return nil, fmt.Errorf("[getUploadID] upload failed code: %v", uploadIDResp.OK)
	}

	// 3. 分片上传
	// 分区
	parts := utils.SplitBytes(content, 10*utils.MB)

	c.logger.Infof("start upload file: %v, parts: %v, size: %.2fMB", filename, len(parts), float64(len(content))/float64(utils.MB))

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
			defer c.logger.Infof("part: %v finished", number)
			errChan <- c.uploadFileClip(preResp.Uri(), preResp.Auth, uploadIDResp.UploadID, number, len(parts), len(part), (number-1)*10*utils.MB, (number-1)*10*utils.MB+len(part), filesize, part)
		}(part, i+1)
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	// 视频上传完成
	checkResp, err := c.uploadCheck(preResp.Uri(), preResp.Auth, filename, uploadIDResp.UploadID, preResp.BizID)
	if err != nil {
		return nil, err
	}

	if checkResp.OK != 1 {
		return nil, fmt.Errorf("[uploadCheck] upload failed code: %v", checkResp.OK)
	}

	c.logger.Infof("video upload finished success，cid: %v", preResp.BizID)

	return &Video{
		Filename: preResp.Filename(),
		Title:    strings.Split(filename, ".")[0],
		Desc:     "",
		CID:      preResp.BizID,
	}, nil
}

/* ===================== helper ===================== */

func (c *Client) getHttpClient(auth bool) *net.HttpClient {
	client := c.httpClient.Clone()

	if c.debug.debug {
		client = client.Debug(c.debug.output)
	}

	if auth && c.authInfo != nil {
		client = client.SetCookies(c.authInfo.Cookies)
	}

	return client
}
