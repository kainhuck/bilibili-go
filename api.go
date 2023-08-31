package bilibili_go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// 获取登陆二维码 https://passport.bilibili.com/x/passport-login/web/qrcode/generate
func (c *Client) qrcodeGenerate() (*QrcodeGenerateResponse, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"

	var baseResp BaseResponse
	err := c.getHttpClient(false).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}

	rsp := &QrcodeGenerateResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// 查询二维码扫描状态 https://passport.bilibili.com/x/passport-login/web/qrcode/poll
func (c *Client) qrcodePoll(qrcodeKey string) (*QrcodePollResponse, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"

	var baseResp BaseResponse
	err := c.getHttpClient(false).Get(uri).AddParams("qrcode_key", qrcodeKey).EndStruct(&baseResp, func(response *http.Response) error {
		c.cookies = response.Cookies()

		return nil
	})
	if err != nil {
		return nil, err
	}

	rsp := &QrcodePollResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetAccount 获取个人账号信息 https://api.bilibili.com/x/member/web/account
func (c *Client) GetAccount() (*AccountResponse, error) {
	uri := "https://api.bilibili.com/x/member/web/account"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &AccountResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetNavigation 获取导航栏信息（个人详细信息） https://api.bilibili.com/x/web-interface/nav
func (c *Client) GetNavigation() (*NavigationResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/nav"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &NavigationResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetNavigationStatus 获取导航栏状态（粉丝数信息）https://api.bilibili.com/x/web-interface/nav/stat
func (c *Client) GetNavigationStatus() (*NavigationStatusResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/nav/stat"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &NavigationStatusResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// 视频预上传 https://member.bilibili.com/preupload
func (c *Client) preUpload(filename string, size int64) (*PreUploadResponse, error) {
	uri := "https://member.bilibili.com/preupload"

	var resp PreUploadResponse

	err := c.getHttpClient(true).Get(uri).
		AddParams("zone", "cs").
		AddParams("upcdn", "bldsa").
		AddParams("probe_version", "20221109").
		AddParams("name", filename).
		AddParams("r", "upos").
		AddParams("profile", "ugcfx/bup").
		AddParams("ssl", "0").
		AddParams("version", "2.14.0.0").
		AddParams("size", strconv.FormatInt(size, 10)).
		AddParams("webVersion", "2.14.0").
		EndStruct(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, err
}

// 获取上传id https://upos-cs-upcdnbldsa.bilivideo.com
func (c *Client) getUploadID(uri string, auth string, bizID int, size int64) (*GetUploadIDResponse, error) {
	var resp GetUploadIDResponse

	err := c.getHttpClient(true).Post(uri).
		SetHeader("X-Upos-Auth", auth).
		AddParams("uploads", "").
		AddParams("output", "json").
		AddParams("profile", "ugcfx/bup").
		AddParams("filesize", strconv.FormatInt(size, 10)).
		AddParams("partsize", "10485760"). // 块大小：10mb
		AddParams("meta_upos_uri", "upos://fxmetalf/n230829qn283p9ffyholy2gigl5advkd.txt").
		AddParams("biz_id", strconv.Itoa(bizID)).
		EndStruct(&resp)

	if err != nil {
		return nil, err
	}

	return &resp, err
}

// 分片上传文件
func (c *Client) uploadFileClip(uri string, auth string, uploadId string, partNumber int, chunks int, size int, start int, end int, total int64, file []byte) error {
	_, _, err := c.getHttpClient(true).Put(uri).
		SetHeader("X-Upos-Auth", auth).
		AddParams("partNumber", strconv.Itoa(partNumber)).
		AddParams("uploadId", uploadId).
		AddParams("chunk", strconv.Itoa(partNumber-1)).
		AddParams("chunks", strconv.Itoa(chunks)).
		AddParams("size", strconv.Itoa(size)).
		AddParams("start", strconv.Itoa(start)).
		AddParams("end", strconv.Itoa(end)).
		AddParams("total", strconv.FormatInt(total, 10)).
		SendBody(bytes.NewReader(file)).
		End()

	return err
}

// 上传完文件后调用该接口
func (c *Client) uploadCheck(uri string, auth string, filename string, uploadID string, bizID int) (*UploadCheckResponse, error) {
	var resp UploadCheckResponse

	err := c.getHttpClient(true).Post(uri).
		SetHeader("X-Upos-Auth", auth).
		AddParams("output", "json").
		AddParams("name", filename).
		AddParams("profile", "ugcfx/bup").
		AddParams("uploadId", uploadID).
		AddParams("biz_id", strconv.Itoa(bizID)).
		EndStruct(&resp)

	if err != nil {
		return nil, err
	}

	return &resp, err
}

// 上传封面 https://member.bilibili.com/x/vu/web/cover/up

// 投稿 https://member.bilibili.com/x/vu/web/add/v3
