package bilibili_go

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// https://passport.bilibili.com/x/passport-login/web/qrcode/generate
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
