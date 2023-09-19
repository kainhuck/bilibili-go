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
	csrf             string
	wbiKey           string // imgKey + subKey
	wbiKeyLastUpdate time.Time
	debug            *debugInfo
	logger           Logger
	showQRCodeFunc   func(code *qrcode.QRCode) error
	mid              int64 // 当前用户mid
	intervalMutex    sync.Mutex
}

func NewClient(opts ...Option) *Client {
	opt := applyOptions(opts...)

	httpClient := net.NewHttpClient(opt.HttpClient).SetUserAgent(opt.UserAgent)

	client := &Client{
		httpClient:     httpClient,
		authStorage:    opt.AuthStorage,
		debug:          opt.Debug,
		logger:         opt.Logger,
		showQRCodeFunc: opt.ShowQRCodeFunc,
		intervalMutex:  sync.Mutex{},
	}

	if opt.RefreshInterval > 0 {
		go func() {
			ticker := time.NewTicker(opt.RefreshInterval)
			for {
				select {
				case <-ticker.C:
					if err := client.RefreshAuthInfo(); err != nil {
						client.logger.Errorf("refresh auth info failed: %v", err)
					}
				}
			}
		}()
	}

	return client
}

func (c *Client) setAuthInfo(auth *AuthInfo) {
	c.authInfo = auth
	if c.authInfo == nil {
		return
	}
	for _, cookie := range c.authInfo.Cookies {
		if cookie.Name == "bili_jct" {
			c.csrf = cookie.Value
		}
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

/* ================= 一下是对接口的二次封装 ================= */

// LoginWithQrCode 登陆这一步必须成功，否则后续接口无法访问
func (c *Client) LoginWithQrCode() {
	if c.authStorage != nil {
		auth, err := c.authStorage.LoadAuthInfo()
		if err == nil && auth != nil {
			c.setAuthInfo(auth)
			user, err := c.GetMyAccount()
			if err == nil {
				c.mid = user.Mid
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
			user, err := c.GetMyAccount()
			if err != nil {
				c.logger.Errorf("login failed: %v", err)
				os.Exit(-1)
			}
			c.mid = user.Mid
			c.logger.Infof("login success!!!")
			return
		case 86038:
			c.logger.Errorf("qrcode expired")
			os.Exit(-1)
		}
		time.Sleep(1 * time.Second)
	}
}

// Logout 退出登陆 会返回重定向链接
func (c *Client) Logout() (string, error) {
	resp, err := c.logout()
	if err != nil {
		return "", err
	}

	if c.authStorage != nil {
		if err := c.authStorage.LogoutAuthInfo(c.authInfo); err != nil {
			c.logger.Errorf("call LogoutAuthInfo failed: %v", err)
		}
	}

	c.authInfo = nil

	return resp.RedirectUrl, nil
}

// UploadVideoFromDisk 从本地磁盘上传视频 videoPath 视频路径
func (c *Client) UploadVideoFromDisk(videoPath string) (*SubmitVideo, error) {
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
func (c *Client) UploadVideoFromReader(filename string, reader io.Reader) (*SubmitVideo, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return c.UploadVideo(filename, content)
}

// UploadVideoFromHTTP 从http链接上传文件
func (c *Client) UploadVideoFromHTTP(filename string, url string) (*SubmitVideo, error) {
	c.logger.Infof("start download file: %v, from: %v", filename, url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.UploadVideoFromReader(filename, resp.Body)
}

// UploadVideo 视频上传，filename 文件名 content 视频内容
func (c *Client) UploadVideo(filename string, content []byte) (*SubmitVideo, error) {
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

	return &SubmitVideo{
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

// Follow 关注用户
func (c *Client) Follow(mid interface{}) error {
	return c.ModifyRelation(mid, 1, 11)
}

// UnFollow 取关用户
func (c *Client) UnFollow(mid interface{}) error {
	return c.ModifyRelation(mid, 2, 11)
}

// WhisperFollow 悄悄关注
func (c *Client) WhisperFollow(mid interface{}) error {
	return c.ModifyRelation(mid, 3, 11)
}

// UnWhisperFollow 取消悄悄关注
func (c *Client) UnWhisperFollow(mid interface{}) error {
	return c.ModifyRelation(mid, 4, 11)
}

// Block 拉黑用户
func (c *Client) Block(mid interface{}) error {
	return c.ModifyRelation(mid, 5, 11)
}

// UnBlock 取消拉黑
func (c *Client) UnBlock(mid interface{}) error {
	return c.ModifyRelation(mid, 6, 11)
}

// GetFollowers 查询自己的粉丝
func (c *Client) GetFollowers(ps int, pn int) (*RelationUserResponse, error) {
	return c.GetUserFollowers(c.mid, ps, pn)
}

// GetFollowings 查询自己的关注
func (c *Client) GetFollowings(orderType string, ps int, pn int) (*RelationUserResponse, error) {
	return c.GetUserFollowings(c.mid, orderType, ps, pn)
}

// GetFollowingsV2 查询自己的关注
func (c *Client) GetFollowingsV2(ps int, pn int) (*RelationUserResponse, error) {
	return c.GetUserFollowingsV2(c.mid, ps, pn)
}

// RefreshAuthInfo 刷新token信息
func (c *Client) RefreshAuthInfo() error {
	if c.authInfo == nil {
		return nil
	}
	c.intervalMutex.Lock()
	defer c.intervalMutex.Unlock()
	// 1. 判断是否需要刷新cookie
	cookieInfo, err := c.getCookieInfo()
	if err != nil {
		return err
	}
	if !cookieInfo.Refresh {
		return nil
	}

	c.logger.Info("refresh auth info")

	oldRefreshToken := c.authInfo.RefreshToken

	// 2. 获取refresh_csrf
	csrf, err := c.getRefreshCSRF()
	if err != nil {
		return err
	}

	// 3. 刷新cookie
	resp, cookies, err := c.refreshCookie(csrf)
	if err != nil {
		return err
	}

	// 4. 保存新的cookie
	c.setAuthInfo(&AuthInfo{
		Cookies:      cookies,
		RefreshToken: resp.RefreshToken,
	})

	// 5. 确认更新
	if err := c.confirmRefresh(oldRefreshToken); err != nil {
		return err
	}

	// 6. 持久化
	if c.authStorage != nil {
		if err := c.authStorage.SaveAuthInfo(c.authInfo); err != nil {
			c.logger.Errorf("SaveAuthInfo failed: %v", err)
		}
	}

	return nil
}

// LikeVideo 点赞视频
func (c *Client) LikeVideo(id string) error {
	return c.likeVideo(id, 1)
}

// UnLikeVideo 取消点赞
func (c *Client) UnLikeVideo(id string) error {
	return c.likeVideo(id, 2)
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
