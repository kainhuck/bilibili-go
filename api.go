package bilibili_go

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// 获取登陆二维码 https://passport.bilibili.com/x/passport-login/web/qrcode/generate
func (c *Client) qrcodeGenerate() (*QrcodeGenerateResponse, error) {
	uri := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate"

	baseResp, err := c.get(uri, nil, nil, nil)
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

	param := make(url.Values)
	param.Add("qrcode_key", qrcodeKey)

	baseResp, err := c.get(uri, param, nil, func(response *http.Response) error {
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

	baseResp, err := c.get(uri, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	rsp := &AccountResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetNavigation 获取导航栏信息（个人详细信息） https://api.bilibili.com/x/web-interface/nav
func (c *Client) GetNavigation() (*NavigationResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/nav"

	baseResp, err := c.get(uri, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	rsp := &NavigationResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// GetNavigationStatus 获取导航栏状态（粉丝数信息）https://api.bilibili.com/x/web-interface/nav/stat
func (c *Client) GetNavigationStatus() (*NavigationStatusResponse, error) {
	uri := "https://api.bilibili.com/x/web-interface/nav/stat"

	baseResp, err := c.get(uri, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	rsp := &NavigationStatusResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}

// 视频预上传 https://member.bilibili.com/preupload
func (c *Client) PreUpload(filename string, size int64) (*NavigationStatusResponse, error) {
	uri := "https://member.bilibili.com/preupload"

	values := url.Values{}
	values.Add("zone", "cs")
	values.Add("upcdn", "bldsa")
	values.Add("probe_version", "20221109")
	values.Add("name", filename)
	values.Add("r", "upos")
	values.Add("profile", "ugcfx/bup")
	values.Add("ssl", "0")
	values.Add("version", "2.14.0.0")
	values.Add("size", strconv.FormatInt(size, 10))
	values.Add("webVersion", "2.14.0")

	var resp PreUploadResponse

	// todo
	baseResp, err := c.get(uri, values, nil, func(response *http.Response) error {
		return json.NewDecoder(response.Body).Decode(&resp)
	})
	if err != nil {
		return nil, err
	}

	rsp := &NavigationStatusResponse{}
	err = json.Unmarshal(baseResp.RawData(), &rsp)

	return rsp, err
}
