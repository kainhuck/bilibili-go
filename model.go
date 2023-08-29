package bilibili_go

import (
	"encoding/json"
	"io"
)

// BaseResponse dor base response
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	TTL     int         `json:"ttl"`
	Data    interface{} `json:"data"`
}

func (r *BaseResponse) RawData() []byte {
	bts, _ := json.Marshal(r.Data)

	return bts
}

func NewBaseResponse(body io.Reader) (*BaseResponse, error) {
	resp := &BaseResponse{}
	err := json.NewDecoder(body).Decode(&resp)

	return resp, err
}

/* ======================================================================= */
/*                          data response                                  */
/* ======================================================================= */

// QrcodeGenerateResponse for qrcode generate response
type QrcodeGenerateResponse struct {
	Url       string `json:"url"`
	QrcodeKey string `json:"qrcode_key"`
}

// QrcodePollResponse for qrcode poll response
type QrcodePollResponse struct {
	Url          string `json:"url"`
	RefreshToken string `json:"refresh_token"`
	Timestamp    int    `json:"timestamp"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}

// AccountResponse for account response
type AccountResponse struct {
	Mid      int64  `json:"mid"`
	UName    string `json:"uname"`
	UserID   string `json:"userid"`
	Sign     string `json:"sign"`
	BirthDay string `json:"birthday"`
	Sex      string `json:"sex"`
	NickFree bool   `json:"nick_free"`
	Rank     string `json:"rank"`
}
