package bilibili_go

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
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
	Uname    string `json:"uname"`
	UserID   string `json:"userid"`
	Sign     string `json:"sign"`
	BirthDay string `json:"birthday"`
	Sex      string `json:"sex"`
	NickFree bool   `json:"nick_free"`
	Rank     string `json:"rank"`
}

// NavigationResponse for navigation response
type NavigationResponse struct {
	AllowanceCount     int            `json:"allowance_count"`
	AnswerStatus       int            `json:"answer_status"`
	EmailVerified      int            `json:"email_verified"`
	Face               string         `json:"face"`
	FaceNFT            int            `json:"face_nft"`
	FaceNFTType        int            `json:"face_nft_type"`
	HasShop            bool           `json:"has_shop"`
	IsLogin            bool           `json:"isLogin"`
	IsJury             bool           `json:"is_jury"`
	IsSeniorMember     int            `json:"is_senior_member"`
	LevelInfo          LevelInfo      `json:"level_info"`
	Mid                int            `json:"mid"`
	MobileVerified     int            `json:"mobile_verified"`
	Money              float64        `json:"money"`
	Moral              int            `json:"moral"`
	Official           Official       `json:"official"`
	OfficialVerify     OfficialVerify `json:"officialVerify"`
	Pendant            Pendant        `json:"pendant"`
	Scores             int            `json:"scores"`
	ShopURL            string         `json:"shop_url"`
	UName              string         `json:"uname"`
	VIP                VIP            `json:"vip"`
	VIPDueDate         int64          `json:"vipDueDate"`
	VIPStatus          int            `json:"vipStatus"`
	VIPType            int            `json:"vipType"`
	VIPAvatarSubscript int            `json:"vip_avatar_subscript"`
	VIPLabel           Label          `json:"vip_label"`
	VIPNicknameColor   string         `json:"vip_nickname_color"`
	VIPPayType         int            `json:"vip_pay_type"`
	VIPThemeType       int            `json:"vip_theme_type"`
	Wallet             Wallet         `json:"wallet"`
	WBIImg             WBIImage       `json:"wbi_img"`
}

type LevelInfo struct {
	CurrentExp   int `json:"current_exp"`
	CurrentLevel int `json:"current_level"`
	CurrentMin   int `json:"current_min"`
	NextExp      int `json:"next_exp"`
}

type Official struct {
	Desc  string `json:"desc"`
	Role  int    `json:"role"`
	Title string `json:"title"`
	Type  int    `json:"type"`
}

type OfficialVerify struct {
	Desc string `json:"desc"`
	Type int    `json:"type"`
}

type Label struct {
	BgColor               string `json:"bg_color"`
	BgStyle               int    `json:"bg_style"`
	BorderColor           string `json:"border_color"`
	ImgLabelUriHans       string `json:"img_label_uri_hans"`
	ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
	ImgLabelUriHant       string `json:"img_label_uri_hant"`
	ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
	LabelTheme            string `json:"label_theme"`
	Path                  string `json:"path"`
	Text                  string `json:"text"`
	TextColor             string `json:"text_color"`
	UseImgLabel           bool   `json:"use_img_label"`
}

type Pendant struct {
	Expire            int    `json:"expire"`
	Image             string `json:"image"`
	ImageEnhance      string `json:"image_enhance"`
	ImageEnhanceFrame string `json:"image_enhance_frame"`
	Name              string `json:"name"`
	PID               int    `json:"pid"`
}

type VIP struct {
	AvatarSubscript    int    `json:"avatar_subscript"`
	AvatarSubscriptURL string `json:"avatar_subscript_url"`
	DueDate            int64  `json:"due_date"`
	Label              Label  `json:"label"`
	NicknameColor      string `json:"nickname_color"`
	Role               int    `json:"role"`
	Status             int    `json:"status"`
	ThemeType          int    `json:"theme_type"`
	TvDueDate          int64  `json:"tv_due_date"`
	TvVIPPayType       int    `json:"tv_vip_pay_type"`
	TvVIPStatus        int    `json:"tv_vip_status"`
	Type               int    `json:"type"`
	VIPPayType         int    `json:"vip_pay_type"`
}

type Wallet struct {
	BCoinBalance  int `json:"bcoin_balance"`
	CouponBalance int `json:"coupon_balance"`
	CouponDueTime int `json:"coupon_due_time"`
	Mid           int `json:"mid"`
}

type WBIImage struct {
	ImgURL string `json:"img_url"`
	SubURL string `json:"sub_url"`
}

// NavigationStatusResponse for navigation status Response
type NavigationStatusResponse struct {
	Following    int64 `json:"following"`     // 关注数
	Follower     int64 `json:"follower"`      // 粉丝数
	DynamicCount int64 `json:"dynamic_count"` // 动态数
}

// PreUploadResponse ...
type PreUploadResponse struct {
	OK              int         `json:"OK"`
	Auth            string      `json:"auth"`
	BizID           int         `json:"biz_id"`
	ChunkRetry      int         `json:"chunk_retry"`
	ChunkRetryDelay int         `json:"chunk_retry_delay"`
	ChunkSize       int         `json:"chunk_size"`
	Endpoint        string      `json:"endpoint"`
	Endpoints       []string    `json:"endpoints"`
	ExposeParams    interface{} `json:"expose_params"`
	PutQuery        string      `json:"put_query"`
	Threads         int         `json:"threads"`
	Timeout         int         `json:"timeout"`
	UIP             string      `json:"uip"`
	UposURI         string      `json:"upos_uri"`
}

func (r *PreUploadResponse) Uri() string {
	return "https:" + r.Endpoint + "/" + strings.TrimPrefix(r.UposURI, "upos://")
}

func (r *PreUploadResponse) Filename() string {
	return strings.Split(filepath.Base(r.UposURI), ".")[0]
}

// GetUploadIDResponse ...
type GetUploadIDResponse struct {
	OK       int    `json:"OK"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	UploadID string `json:"upload_id"`
}

// UploadCheckResponse ...
type UploadCheckResponse struct {
	OK       int    `json:"OK"`
	Bucket   string `json:"bucket"`
	Etag     string `json:"etag"`
	Key      string `json:"key"`
	Location string `json:"location"`
}

// UploadCoverResponse ...
type UploadCoverResponse struct {
	Url string `json:"url"`
}

type Video struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	CID      int    `json:"cid"`
}

type Subtitle struct {
	Open int    `json:"open"`
	Lan  string `json:"lan"`
}

// SubmitRequest ...
type SubmitRequest struct {
	Cover            string   `json:"cover"`              // 封面 必须
	Title            string   `json:"title"`              // 标题 必须
	Copyright        int      `json:"copyright"`          // 是否原创 必须
	TID              int      `json:"tid"`                // 分类ID 必须
	Tag              string   `json:"tag"`                // 标签 用逗号分隔 必须
	DescFormatID     int      `json:"desc_format_id"`     // ？
	Desc             string   `json:"desc"`               // 简介 必须
	Recreate         int      `json:"recreate"`           // 二创视频 ？
	Dynamic          string   `json:"dynamic"`            // 粉丝动态 ？
	Interactive      int      `json:"interactive"`        // 是否是合作视频 ？
	Videos           []*Video `json:"videos"`             // 视频 必须
	ActReserveCreate int      `json:"act_reserve_create"` // 允许二创 ？
	NoDisturbance    int      `json:"no_disturbance"`     // ？
	NoReprint        int      `json:"no_reprint"`         // ？
	Subtitle         Subtitle `json:"subtitle"`           // ？
	Dolby            int      `json:"dolby"`              // 杜比音效
	LosslessMusic    int      `json:"lossless_music"`     // 无损音质 ？
	WebOS            int      `json:"web_os"`             // ？2
	CSRF             string   `json:"csrf"`               // bili_jct
}

// SubmitResponse ...
type SubmitResponse struct {
	Aid  int64  `json:"aid"`
	Bvid string `json:"bvid"`
}
