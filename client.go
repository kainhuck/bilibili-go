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
	authInfo         *AuthInfo
	authStorage      AuthStorage
	cookieCache      map[string]string
	wbiKey           string // imgKey + subKey
	wbiKeyLastUpdate time.Time
	debug            *debugInfo
	logger           Logger
	showQRCodeFunc   func(code *qrcode.QRCode) error
}

func NewClient(opts ...Option) *Client {
	opt := applyOptions(opts...)

	client := net.NewHttpClient(opt.HttpClient).SetUserAgent(opt.UserAgent)

	return &Client{
		httpClient:     client,
		cookieCache:    make(map[string]string),
		authStorage:    opt.AuthStorage,
		debug:          opt.Debug,
		logger:         opt.Logger,
		showQRCodeFunc: opt.ShowQRCodeFunc,
	}
}

func (c *Client) setAuthInfo(auth *AuthInfo) {
	c.authInfo = auth
	if c.authInfo == nil {
		return
	}
	for _, cookie := range c.authInfo.Cookies {
		c.cookieCache[cookie.Name] = cookie.Value
	}
}

func (c *Client) getWbiKey() string {
	resp, err := c.GetNavigation()
	if err != nil {
		return ""
	}

	imgKey := strings.Split(strings.Split(resp.WBIImg.ImgURL, "/")[len(strings.Split(resp.WBIImg.ImgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(resp.WBIImg.SubURL, "/")[len(strings.Split(resp.WBIImg.SubURL, "/"))-1], ".")[0]
	return imgKey + subKey
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
	if c.authStorage != nil {
		auth, err := c.authStorage.LoadAuthInfo()
		if err == nil && auth != nil {
			c.setAuthInfo(auth)
			user, err := c.getMyInfo()
			if err == nil {
				c.authInfo.User = user
				if err := c.authStorage.SaveAuthInfo(c.authInfo); err != nil {
					c.logger.Errorf("SaveAuthInfo failed: %v", err)
					os.Exit(-1)
				}
				c.logger.Info("load auth info from storage")
				return
			} else {
				// maybe token过期
				c.logger.Warnf("auth info error: %v", err)
			}
		}
		if err != nil {
			c.logger.Errorf("load auth info failed: %v", err)
		}
	}

	defer func() {
		if c.authStorage != nil {
			if err := c.authStorage.SaveAuthInfo(c.authInfo); err != nil {
				c.logger.Errorf("SaveAuthInfo failed: %v", err)
			}
		}
	}()

	generateResp, err := c.qrcodeGenerate()
	if err != nil {
		c.logger.Errorf("generate qrcode failed: %v", err)
		os.Exit(-1)
	}

	qrCode, err := qrcode.New(generateResp.Url, qrcode.Medium)
	if err != nil {
		c.logger.Errorf("new qrcode failed: %v", err)
		os.Exit(-1)
	}

	if err := c.showQRCodeFunc(qrCode); err != nil {
		c.logger.Errorf("show qrcode failed: %v", err)
		os.Exit(-1)
	}

	for {
		resp, cookies, err := c.qrcodePoll(generateResp.QrcodeKey)
		if err != nil {
			c.logger.Errorf("poll qrcode failed: %v", err)
			os.Exit(-1)
		}

		switch resp.Code {
		case 0:
			c.setAuthInfo(&AuthInfo{
				Cookies:      cookies,
				RefreshToken: resp.RefreshToken,
			})
			resp, err := c.getMyInfo()
			if err != nil {
				c.logger.Errorf("login failed: %v", err)
				os.Exit(-1)
			}
			c.authInfo.User = resp
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

// UploadCoverFromDisk 从本地磁盘上传封面 imagePath 图片路径
func (c *Client) UploadCoverFromDisk(imagePath string) (*UploadCoverResponse, error) {
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, err
	}

	return c.UploadCover(imageData)
}

// UploadCoverFromReader ...
func (c *Client) UploadCoverFromReader(reader io.Reader) (*UploadCoverResponse, error) {
	imageData, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return c.UploadCover(imageData)
}

// UploadCoverFromHTTP 从http链接上传封面
func (c *Client) UploadCoverFromHTTP(url string) (*UploadCoverResponse, error) {
	c.logger.Infof("start download cover from: %v", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.UploadCoverFromReader(resp.Body)
}

// GetMyInfo 获取当前用户信息
func (c *Client) GetMyInfo() (*GetMyInfoResponse, error) {
	if c.authInfo.User != nil {
		return c.authInfo.User, nil
	}
	user, err := c.getMyInfo()
	if err != nil {
		return nil, err
	}
	c.authInfo.User = user

	_ = c.authStorage.SaveAuthInfo(c.authInfo)

	return user, nil
}

// Follow 关注用户
func (c *Client) Follow(mid string) error {
	return c.ModifyRelation(mid, 1, 11)
}

// UnFollow 取关用户
func (c *Client) UnFollow(mid string) error {
	return c.ModifyRelation(mid, 2, 11)
}

// WhisperFollow 悄悄关注
func (c *Client) WhisperFollow(mid string) error {
	return c.ModifyRelation(mid, 3, 11)
}

// UnWhisperFollow 取消悄悄关注
func (c *Client) UnWhisperFollow(mid string) error {
	return c.ModifyRelation(mid, 4, 11)
}

// Block 拉黑用户
func (c *Client) Block(mid string) error {
	return c.ModifyRelation(mid, 5, 11)
}

// UnBlock 取消拉黑
func (c *Client) UnBlock(mid string) error {
	return c.ModifyRelation(mid, 6, 11)
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
