package bilibili_go

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
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
	err := c.getHttpClient(false).Get(uri).
		AddParams("qrcode_key", qrcodeKey).
		EndStruct(&baseResp, func(response *http.Response) error {
			c.setCookies(response.Cookies())

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

// UploadCover 上传封面 https://member.bilibili.com/x/vu/web/cover/up
func (c *Client) UploadCover(image string) (*UploadCoverResponse, error) {
	uri := "https://member.bilibili.com/x/vu/web/cover/up"

	imageData, err := os.ReadFile(image)
	if err != nil {
		return nil, err
	}
	base64Str := base64.StdEncoding.EncodeToString(imageData)

	var baseResp BaseResponse

	err = c.getHttpClient(true).Post(uri).
		AddParams("t", strconv.FormatInt(time.Now().UnixMilli(), 10)).
		AddFormData("cover", "data:image/jpeg;base64,"+base64Str).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &UploadCoverResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// SubmitVideo 视频投稿 https://member.bilibili.com/x/vu/web/add/v3
func (c *Client) SubmitVideo(req *SubmitRequest) (*SubmitResponse, error) {
	uri := "https://member.bilibili.com/x/vu/web/add/v3"

	req.CSRF = c.cookieCache["bili_jct"]

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	var baseResp BaseResponse

	err = c.getHttpClient(true).
		SetContentType("application/json;charset=UTF-8").
		Post(uri).
		AddParams("t", strconv.FormatInt(time.Now().UnixMilli(), 10)).
		AddParams("csrf", c.cookieCache["bili_jct"]).
		SendBody(bytes.NewReader(reqData)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}

	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &SubmitResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetCoin 获取硬币数 https://account.bilibili.com/site/getCoin
func (c *Client) GetCoin() (*GetCoinResponse, error) {
	uri := "https://account.bilibili.com/site/getCoin"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetCoinResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetSpaceInfo 用户空间详细信息 https://api.bilibili.com/x/space/wbi/acc/info
func (c *Client) GetSpaceInfo(mid string) (*GetSpaceInfoResponse, error) {
	uri := "https://api.bilibili.com/x/space/wbi/acc/info"

	params := make(url.Values)
	params.Add("mid", mid)
	encWbi(params, c.getWbiKeyCached())

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).CoverParams(params).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetSpaceInfoResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserCard 用户名片信息 https://api.bilibili.com/x/web-interface/card
//
//	mid 用户mid
//	photo 是否请求用户主页头像
func (c *Client) GetUserCard(mid string, photo bool) (*GetUserCardResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/card"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("mid", mid).
		AddParams("photo", strconv.FormatBool(photo)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetUserCardResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetMyInfo 登陆用户空间详细信息 https://api.bilibili.com/x/space/myinfo
func (c *Client) GetMyInfo() (*GetMyInfoResponse, error) {
	uri := "https://api.bilibili.com/x/space/myinfo"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetMyInfoResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}
