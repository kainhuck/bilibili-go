package bilibili_go

import (
	"encoding/json"
	"net/http"
	"net/url"
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

	baseResp, err := c.get(uri, param, nil, func(response *http.Response) {
		c.cookies = response.Cookies()
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

// GetNavigation 获取导航栏信息 https://api.bilibili.com/x/web-interface/nav
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
