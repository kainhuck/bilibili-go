package bilibili_go

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kainhuck/bilibili-go/internal/utils"
	"github.com/spf13/cast"
	"net/http"
	"strconv"
	"strings"
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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

	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &GetCoinResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserInfo 用户空间详细信息 https://api.bilibili.com/x/space/wbi/acc/info
func (c *Client) GetUserInfo(mid interface{}) (*GetUserInfoResponse, error) {
	uri := "https://api.bilibili.com/x/space/wbi/acc/info"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		SetWbiKey(c.getWbiKeyCached()).
		AddParams("mid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &GetUserInfoResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUserCard 用户名片信息 https://api.bilibili.com/x/web-interface/card
//
//	mid 用户mid
//	photo 是否请求用户主页头像
func (c *Client) GetUserCard(mid interface{}, photo bool) (*GetUserCardResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/card"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("mid", cast.ToString(mid)).
		AddParams("photo", strconv.FormatBool(photo)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &GetMyInfoResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetRelationStat 获取用户关系状态 https://api.bilibili.com/x/relation/stat
func (c *Client) GetRelationStat(mid interface{}) (*GetRelationStatResponse, error) {
	uri := "https://api.bilibili.com/x/relation/stat"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &GetRelationStatResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetUpStat 获取up主状态数 https://api.bilibili.com/x/space/upstat
func (c *Client) GetUpStat(mid interface{}) (*GetUpStatResponse, error) {
	uri := "https://api.bilibili.com/x/space/upstat"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("mid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &GetUpStatResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetDocUploadCount 相簿投稿数 https://api.vc.bilibili.com/link_draw/v1/doc/upload_count
func (c *Client) GetDocUploadCount(mid interface{}) (*GetDocUploadCountResponse, error) {
	uri := "https://api.vc.bilibili.com/link_draw/v1/doc/upload_count"

	var baseResp BaseResponse
	err := c.getHttpClient(false).Get(uri).
		AddParams("uid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
func (c *Client) GetUserFollowers(mid interface{}, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followers"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
func (c *Client) GetUserFollowings(mid interface{}, orderType string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followings"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		AddParams("order_type", orderType).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
func (c *Client) GetUserFollowingsV2(mid interface{}, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://app.biliapi.net/x/v2/relation/followings"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
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
func (c *Client) SearchUserFollowings(mid interface{}, name string, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/followings/search"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		AddParams("name", name).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetSameFollowings 查询共同关注列表 https://api.bilibili.com/x/relation/same/followings
// mid 目标用户ID
// ps 每页大小
// pn 页码
func (c *Client) GetSameFollowings(mid interface{}, ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/same/followings"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("vmid", cast.ToString(mid)).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetWhispers 查询悄悄关注列表 https://api.bilibili.com/x/relation/whispers
// mid 目标用户ID
// ps 每页大小
// pn 页码
// 只能查看自己的悄悄关注，total字段不返回，list 返回全部
func (c *Client) GetWhispers() (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/whispers"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetFriends 查询互相关注列表 https://api.bilibili.com/x/relation/friends
// mid 目标用户ID
// ps 每页大小
// pn 页码
// 只能查看自己的互相关注，total字段不返回，list 返回全部
func (c *Client) GetFriends() (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/friends"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetBlacks 查询黑名单列表 https://api.bilibili.com/x/relation/blacks
// ps 每页大小
// pn 页码
func (c *Client) GetBlacks(ps int, pn int) (*RelationUserResponse, error) {
	uri := "https://api.bilibili.com/x/relation/blacks"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &RelationUserResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// ModifyRelation 操作用户关系 https://api.bilibili.com/x/relation/modify
// mid 目标用户mid
// act 操作代码
//
//	1 关注
//	2 取关
//	3 悄悄关注
//	4 取消悄悄关注
//	5 拉黑
//	6 取消拉黑
//	7 踢出粉丝
//
// reSrc 关注来源
//
//	11 空间
//	14 视频
//	115 文章
//	222 活动页面
func (c *Client) ModifyRelation(mid interface{}, act int, reSrc int) error {
	uri := "https://api.bilibili.com/x/relation/modify"

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("fid", cast.ToString(mid)).
		AddFormData("act", strconv.Itoa(act)).
		AddFormData("re_src", strconv.Itoa(reSrc)).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// BatchModifyRelation 批量操作用户关系 https://api.bilibili.com/x/relation/batch/modify
// mids 目标用户mid
// act 操作代码
//
//	1 关注
//	5 拉黑
//
// reSrc 关注来源
//
//	11 空间
//	14 视频
//	115 文章
//	222 活动页面
func (c *Client) BatchModifyRelation(mids []string, act int, reSrc int) (*BatchModifyRelationResponse, error) {
	uri := "https://api.bilibili.com/x/relation/batch/modify"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Post(uri).
		AddFormData("fids", strings.Join(mids, ",")).
		AddFormData("act", strconv.Itoa(act)).
		AddFormData("re_src", strconv.Itoa(reSrc)).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &BatchModifyRelationResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetRelation 查询用户与自己的关系 https://api.bilibili.com/x/relation
// mid 用户ID
func (c *Client) GetRelation(mid interface{}) (*Relation, error) {
	uri := "https://api.bilibili.com/x/relation"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("fid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &Relation{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetAccRelation 查询用户与自己的互相关系 https://api.bilibili.com/x/space/wbi/acc/relation
func (c *Client) GetAccRelation(mid interface{}) (*AccRelation, error) {
	uri := "https://api.bilibili.com/x/space/wbi/acc/relation"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("mid", cast.ToString(mid)).
		SetWbiKey(c.getWbiKeyCached()).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &AccRelation{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// BatchGetRelation 批量查询用户与自己的关系 https://api.bilibili.com/x/relation/relations
// 返回的key是mid
func (c *Client) BatchGetRelation(mid ...string) (map[string]Relation, error) {
	uri := "https://api.bilibili.com/x/relation/relations"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("fids", strings.Join(mid, ",")).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := make(map[string]Relation)
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetRelationTags 查询关注分组列表 https://api.bilibili.com/x/relation/tags
func (c *Client) GetRelationTags() ([]*RelationTag, error) {
	uri := "https://api.bilibili.com/x/relation/tags"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := make([]*RelationTag, 0)
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetRelationTagUsers 查询关注分组内的用户 https://api.bilibili.com/x/relation/tag
// tagId 关注分组id 可通过 GetRelationTags 接口获取
// orderType 按照关注顺序排列：留空 按照最常访问排列：attention
// ps 每页项数
// pn 页码
func (c *Client) GetRelationTagUsers(tagId int, orderType string, ps int, pn int) ([]*RelationUser, error) {
	uri := "https://api.bilibili.com/x/relation/tag"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("tagid", strconv.Itoa(tagId)).
		AddParams("order_type", orderType).
		AddParams("ps", strconv.Itoa(ps)).
		AddParams("pn", strconv.Itoa(pn)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := make([]*RelationUser, 0)
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// QueryRelationTagByUser 查询用户所在的分组 https://api.bilibili.com/x/relation/tag/user
// mid 用户ID
// 返回的 key 是分组ID， value 是分组名称
func (c *Client) QueryRelationTagByUser(mid interface{}) (map[string]string, error) {
	uri := "https://api.bilibili.com/x/relation/tag/user"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("fid", cast.ToString(mid)).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := make(map[string]string)
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetSpecialRelationTagUsers 查询特别关注的所有用户mid https://api.bilibili.com/x/relation/tag/special
// 返回所有用户的mid
func (c *Client) GetSpecialRelationTagUsers() ([]string, error) {
	uri := "https://api.bilibili.com/x/relation/tag/special"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := make([]string, 0)
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// CreateRelationTag 创建分组 https://api.bilibili.com/x/relation/tag/create
// name 分组名称
func (c *Client) CreateRelationTag(name string) (*CreateRelationTagResponse, error) {
	uri := "https://api.bilibili.com/x/relation/tag/create"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Post(uri).
		AddFormData("tag", name).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &CreateRelationTagResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// UpdateRelationTag 更新分组 https://api.bilibili.com/x/relation/tag/update
// tagId 分组ID
// name 分组新名称
func (c *Client) UpdateRelationTag(tagId int, name string) error {
	uri := "https://api.bilibili.com/x/relation/tag/update"

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("tagid", strconv.Itoa(tagId)).
		AddFormData("name", name).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// DeleteRelationTag 删除分组 https://api.bilibili.com/x/relation/tag/del
func (c *Client) DeleteRelationTag(tagId int) error {
	uri := "https://api.bilibili.com/x/relation/tag/del"

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("tagid", strconv.Itoa(tagId)).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// AddUsersToRelationTags 往分组内添加成员 https://api.bilibili.com/x/relation/tags/addUsers
// 通过该接口可以将多个用户移动到多个分组
// 如需移除分组中的成员，请将tagids设为 0，即移动至默认分组，而不是取关
// mids 用户ID
// tagIds 分组ID
func (c *Client) AddUsersToRelationTags(mids []string, tagIds []int) error {
	uri := "https://api.bilibili.com/x/relation/tags/addUsers"

	tagIdsString := make([]string, 0, len(tagIds))
	for _, each := range tagIds {
		tagIdsString = append(tagIdsString, strconv.Itoa(each))
	}

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("fids", strings.Join(mids, ",")).
		AddFormData("tagids", strings.Join(tagIdsString, ",")).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// CopyUsersToRelationTags 复制成员到分组 https://api.bilibili.com/x/relation/tags/copyUsers
// mids 用户ID
// tagIds 分组ID
func (c *Client) CopyUsersToRelationTags(mids []string, tagIds []int) error {
	uri := "https://api.bilibili.com/x/relation/tags/copyUsers"

	tagIdsString := make([]string, 0, len(tagIds))
	for _, each := range tagIds {
		tagIdsString = append(tagIdsString, strconv.Itoa(each))
	}

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("fids", strings.Join(mids, ",")).
		AddFormData("tagids", strings.Join(tagIdsString, ",")).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// MoveUsersToRelationTags 复制成员到分组 https://api.bilibili.com/x/relation/tags/moveUsers
// mids 用户ID
// tagIds 分组ID
func (c *Client) MoveUsersToRelationTags(mids []string, beforeTagIds []int, afterTagIds []int) error {
	uri := "https://api.bilibili.com/x/relation/tags/moveUsers"

	beforeTagIdsString := make([]string, 0, len(beforeTagIds))
	for _, each := range beforeTagIds {
		beforeTagIdsString = append(beforeTagIdsString, strconv.Itoa(each))
	}

	afterTagIdsString := make([]string, 0, len(afterTagIds))
	for _, each := range afterTagIds {
		afterTagIdsString = append(afterTagIdsString, strconv.Itoa(each))
	}

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("fids", strings.Join(mids, ",")).
		AddFormData("beforeTagids", strings.Join(beforeTagIdsString, ",")).
		AddFormData("afterTagids", strings.Join(afterTagIdsString, ",")).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}

// logout 退出登陆 https://passport.bilibili.com/login/exit/v2
func (c *Client) logout() (*LogoutResponse, error) {
	uri := "https://passport.bilibili.com/login/exit/v2"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Post(uri).
		AddFormData("biliCSRF", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &LogoutResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// getCookieInfo 检查是否需要刷新cookie https://passport.bilibili.com/x/passport-login/web/cookie/info
func (c *Client) getCookieInfo() (*CookieInfo, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/cookie/info"

	var baseResp BaseResponse
	err := c.getHttpClient(true).Get(uri).
		AddParams("biliCSRF", c.cookieCache["bili_jct"]).
		EndStruct(&baseResp)
	if err != nil {
		return nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, fmt.Errorf("%s", bts)
	}

	rsp := &CookieInfo{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// getRefreshCSRF 获取 refresh_csrf
func (c *Client) getRefreshCSRF() (string, error) {
	path, err := utils.GetCorrespondPath(time.Now().UnixMilli())
	if err != nil {
		return "", err
	}

	uri := "https://www.bilibili.com/correspond/1/" + path

	_, body, err := c.getHttpClient(true).Get(uri).End()
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	return doc.Find("#1-name").First().Text(), err
}

// refreshCookie 刷新cookie
func (c *Client) refreshCookie(refreshCsrf string) (*RefreshCookieResponse, []*http.Cookie, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/cookie/refresh"

	var baseResp BaseResponse
	var cookies []*http.Cookie

	err := c.getHttpClient(true).Post(uri).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		AddFormData("refresh_csrf", refreshCsrf).
		AddFormData("refresh_token", c.authInfo.RefreshToken).
		EndStruct(&baseResp, func(response *http.Response) error {
			cookies = response.Cookies()

			return nil
		})
	if err != nil {
		return nil, nil, err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return nil, nil, fmt.Errorf("%s", bts)
	}

	rsp := &RefreshCookieResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, cookies, err

}

// confirmRefresh 确认更新
func (c *Client) confirmRefresh(refreshToken string) error {
	uri := "https://passport.bilibili.com/x/passport-login/web/confirm/refresh"

	var baseResp BaseResponse

	err := c.getHttpClient(true).Post(uri).
		AddFormData("csrf", c.cookieCache["bili_jct"]).
		AddFormData("refresh_token", refreshToken).
		EndStruct(&baseResp)
	if err != nil {
		return err
	}
	if baseResp.Code != CodeSuccess {
		bts, _ := json.Marshal(baseResp)
		return fmt.Errorf("%s", bts)
	}

	return nil
}
