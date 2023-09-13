package bilibili_go

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
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
func (c *Client) qrcodePoll(qrcodeKey string) (*QrcodePollResponse, []*http.Cookie, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/poll"

	var baseResp BaseResponse
	var cookies []*http.Cookie

	err := c.getHttpClient(false).Get(uri).
		AddParams("qrcode_key", qrcodeKey).
		EndStruct(&baseResp, func(response *http.Response) error {
			cookies = response.Cookies()

			return nil
		})
	if err != nil {
		return nil, nil, err
	}

	rsp := &QrcodePollResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, cookies, err
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
func (c *Client) UploadCover(imageData []byte) (*UploadCoverResponse, error) {
	uri := "https://member.bilibili.com/x/vu/web/cover/up"

	base64Str := base64.StdEncoding.EncodeToString(imageData)

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
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

// GetUserInfo 用户空间详细信息 https://api.bilibili.com/x/space/wbi/acc/info
func (c *Client) GetUserInfo(mid string) (*GetUserInfoResponse, error) {
	uri := "https://api.bilibili.com/x/space/wbi/acc/info"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		SetWbiKey(c.getWbiKeyCached()).
		AddParams("mid", mid).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetUserInfoResponse{}
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

// getMyInfo 登陆用户空间详细信息 https://api.bilibili.com/x/space/myinfo
func (c *Client) getMyInfo() (*GetMyInfoResponse, error) {
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

// GetRelationStat 获取用户关系状态 https://api.bilibili.com/x/relation/stat
func (c *Client) GetRelationStat(mid string) (*GetRelationStatResponse, error) {
	uri := "https://api.bilibili.com/x/relation/stat"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", mid).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetRelationStatResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUpStat 获取up主状态数 https://api.bilibili.com/x/space/upstat
func (c *Client) GetUpStat(mid string) (*GetUpStatResponse, error) {
	uri := "https://api.bilibili.com/x/space/upstat"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("mid", mid).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetUpStatResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetDocUploadCount 相簿投稿数 https://api.vc.bilibili.com/link_draw/v1/doc/upload_count
func (c *Client) GetDocUploadCount(mid string) (*GetDocUploadCountResponse, error) {
	uri := "https://api.vc.bilibili.com/link_draw/v1/doc/upload_count"

	var baseResp BaseResponse
	err := c.getHttpClient(false).Get(uri).
		AddParams("uid", mid).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &GetDocUploadCountResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserFollowers 查询用户粉丝列表 https://api.bilibili.com/x/relation/followers
// mid 用户ID
// ps 每页大小
// pn 页码
// 注意：查询别的用户粉丝数上限为250
func (c *Client) GetUserFollowers(mid string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followers"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", mid).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserFollowings 查询用户关注列表 https://api.bilibili.com/x/relation/followings
// mid 用户ID
// orderType 排序方式  按照关注顺序排列：留空  按照最常访问排列：attention
// ps 每页大小
// pn 页码
// 注意：查询别的用户关注数上限为250
func (c *Client) GetUserFollowings(mid string, orderType string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followings"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", mid).
		AddParams("order_type", orderType).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserFollowingsV2 查询用户关注列表 https://app.biliapi.net/x/v2/relation/followings
// mid 用户ID
// ps 每页大小
// pn 页码
// 注意：仅可查看前 5 页 可以获取已设置可见性隐私的关注列表
func (c *Client) GetUserFollowingsV2(mid string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://app.biliapi.net/x/v2/relation/followings"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", mid).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// SearchUserFollowings 搜索用户关注列表 https://api.bilibili.com/x/relation/followings/search
// mid 目标用户ID
// name 搜索关键词
// ps 每页大小
// pn 页码
func (c *Client) SearchUserFollowings(mid string, name string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followings/search"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", mid).
		AddParams("name", name).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != 0 {
		return nil, fmt.Errorf(baseResp.Message)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}
