package bilibili_go

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

type Code int

const (
	// CodeSuccess 成功
	CodeSuccess Code = 0
	// CodeCsrfFailed csrf校验失败
	CodeCsrfFailed Code = -111
	// CodeUnLogin 账号未登录
	CodeUnLogin Code = -101
	// CodeRequestError 请求错误
	CodeRequestError Code = -400
	// CodePermissionDenied 没有权限
	CodePermissionDenied Code = 22104
	// CodeUnFollowed 未关注
	CodeUnFollowed Code = 22105
)

// BaseResponse dor base response
type BaseResponse struct {
	Code    Code        `json:"code"`
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
	AllowanceCount int    `json:"allowance_count"`  // ？
	AnswerStatus   int    `json:"answer_status"`    // ？
	EmailVerified  int    `json:"email_verified"`   // 是否验证邮箱地址，0 未验证 1 已验证
	Face           string `json:"face"`             // 头像
	FaceNFT        int    `json:"face_nft"`         // 是否为NFT头像 0 不是 1 是
	FaceNFTType    int    `json:"face_nft_type"`    // NFT头像类型？
	HasShop        bool   `json:"has_shop"`         // 是否拥有推广商品 true 有 false 无
	IsLogin        bool   `json:"isLogin"`          // 是否已登陆 true 已登陆 false 未登录
	IsJury         bool   `json:"is_jury"`          // 是否是风纪委员 true 是 false 不是
	IsSeniorMember int    `json:"is_senior_member"` // 是否是硬核会员 0 不是 1 是
	LevelInfo      struct {
		CurrentExp   int `json:"current_exp"`
		CurrentLevel int `json:"current_level"`
		CurrentMin   int `json:"current_min"`
		NextExp      int `json:"next_exp"`
	} `json:"level_info"` // 等级信息
	Mid            int64   `json:"mid"`             // 用户 mid
	MobileVerified int     `json:"mobile_verified"` // 是否验证手机号 0 未验证 1 已验证
	Money          float64 `json:"money"`           // 硬币数
	Moral          int     `json:"moral"`           // 当前节操值 上限70
	Official       struct {
		Desc  string `json:"desc"`
		Role  int    `json:"role"`
		Title string `json:"title"`
		Type  int    `json:"type"`
	} `json:"official"` // 认证信息
	OfficialVerify struct {
		Desc string `json:"desc"`
		Type int    `json:"type"`
	} `json:"officialVerify"` // 认证信息2
	Pendant struct {
		Expire            int    `json:"expire"`
		Image             string `json:"image"`
		ImageEnhance      string `json:"image_enhance"`
		ImageEnhanceFrame string `json:"image_enhance_frame"`
		Name              string `json:"name"`
		PID               int    `json:"pid"`
	} `json:"pendant"` // 头像框信息
	Scores  int    `json:"scores"`   // ？
	ShopURL string `json:"shop_url"` // 商品推广页url
	UName   string `json:"uname"`    // 用户昵称
	VIP     struct {
		AvatarSubscript    int    `json:"avatar_subscript"`
		AvatarSubscriptURL string `json:"avatar_subscript_url"`
		DueDate            int64  `json:"due_date"`
		Label              struct {
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
		} `json:"label"`
		NicknameColor string `json:"nickname_color"`
		Role          int    `json:"role"`
		Status        int    `json:"status"`
		ThemeType     int    `json:"theme_type"`
		TvDueDate     int64  `json:"tv_due_date"`
		TvVIPPayType  int    `json:"tv_vip_pay_type"`
		TvVIPStatus   int    `json:"tv_vip_status"`
		Type          int    `json:"type"`
		VIPPayType    int    `json:"vip_pay_type"`
	} `json:"vip"` // 会员信息
	VIPDueDate         int64 `json:"vipDueDate"`           // 会员到期时间 毫秒时间戳
	VIPStatus          int   `json:"vipStatus"`            // 会员开通状态 0 无 1 有
	VIPType            int   `json:"vipType"`              // 会员类型 0 无 1 月度大会员 2 年度及以上大会员
	VIPAvatarSubscript int   `json:"vip_avatar_subscript"` // 是否显示会员图标 0 不显示 1 显示
	VIPLabel           struct {
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
	} `json:"vip_label"` // 会员标签
	VIPNicknameColor string `json:"vip_nickname_color"` // 会员昵称颜色 颜色码
	VIPPayType       int    `json:"vip_pay_type"`       // 会员开通状态 0 无 1 有
	VIPThemeType     int    `json:"vip_theme_type"`     // ？
	Wallet           struct {
		BCoinBalance  int   `json:"bcoin_balance"`
		CouponBalance int   `json:"coupon_balance"`
		CouponDueTime int   `json:"coupon_due_time"`
		Mid           int64 `json:"mid"`
	} `json:"wallet"` // B币钱包信息
	WBIImg struct {
		ImgURL string `json:"img_url"`
		SubURL string `json:"sub_url"`
	} `json:"wbi_img"` // Wbi签名实时口令
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

// SubmitVideo 投稿视频
type SubmitVideo struct {
	Filename string `json:"filename"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	CID      int    `json:"cid"`
}

// SubmitRequest ...
type SubmitRequest struct {
	Cover            string         `json:"cover"`              // 封面 必须
	Title            string         `json:"title"`              // 标题 必须
	Copyright        int            `json:"copyright"`          // 是否原创 必须 1 原创 2 转载
	Source           string         `json:"source"`             // 如果选择转载则 将原视频链接贴这里
	Dtime            int64          `json:"dtime"`              // 设置定时发布时间，时间戳，精确到秒，如不设置则立即发送
	TID              int            `json:"tid"`                // 分类ID 必须
	Tag              string         `json:"tag"`                // 标签 用逗号分隔 必须
	DescFormatID     int            `json:"desc_format_id"`     // ？
	Desc             string         `json:"desc"`               // 简介 必须
	Recreate         int            `json:"recreate"`           // 二创视频 -1 不允许二创 1 允许二创
	Dynamic          string         `json:"dynamic"`            // 粉丝动态 ？
	Interactive      int            `json:"interactive"`        // 是否是合作视频 ？
	Videos           []*SubmitVideo `json:"videos"`             // 视频 必须
	ActReserveCreate int            `json:"act_reserve_create"` // 允许二创 ？
	NoDisturbance    int            `json:"no_disturbance"`     // ？
	NoReprint        int            `json:"no_reprint"`         // ？
	Subtitle         struct {
		Open int    `json:"open"`
		Lan  string `json:"lan"`
	} `json:"subtitle"` // ？
	Dolby         int    `json:"dolby"`          // 杜比音效
	LosslessMusic int    `json:"lossless_music"` // 无损音质 ？
	WebOS         int    `json:"web_os"`         // ？2
	CSRF          string `json:"csrf"`           // bili_jct
}

// SubmitResponse ...
type SubmitResponse struct {
	Aid  int64  `json:"aid"`
	Bvid string `json:"bvid"`
}

// GetCoinResponse 用户硬币信息
type GetCoinResponse struct {
	Money float64 `json:"money"`
}

// GetUserInfoResponse 用户空间信息
type GetUserInfoResponse struct {
	Mid         int64    `json:"mid"`           // mid
	Name        string   `json:"name"`          // 昵称
	Sex         string   `json:"sex"`           // 性别 男 女 保密
	Face        string   `json:"face"`          // 头像链接
	FaceNft     int      `json:"face_nft"`      // 是否为NFT头像 0 不是 1 是
	FaceNftType int      `json:"face_nft_type"` // nft头像类型?
	Sign        string   `json:"sign"`          // 签名
	Rank        int      `json:"rank"`          // 用户权限等级 5000 0级未答题 10000 普通会员 20000 字幕君 25000 VIP 30000 真.职人 32000管理员
	Level       int      `json:"level"`         // 当前等级 0-6 级
	JoinTime    int      `json:"jointime"`      // 注册时间 此接口返回恒为0
	Moral       int      `json:"moral"`         // 节操值 此接口返回恒为0
	Silence     int      `json:"silence"`       // 封禁状态 0 正常 1 被封
	Coins       float64  `json:"coins"`         // 硬币数 需要登陆，只能查看自己的，默认为0
	FansBadge   bool     `json:"fans_badge"`    // 是否具有粉丝勋章 false 无 true 有
	FansMedal   struct { // 粉丝勋章信息
		Show  bool        `json:"show"`
		Wear  bool        `json:"wear"`
		Medal interface{} `json:"medal"`
	} `json:"fans_medal"`
	Official struct { // 认证信息
		Role  int    `json:"role"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Type  int    `json:"type"`
	} `json:"official"`
	VIP struct { // 会员信息
		Type      int   `json:"type"`
		Status    int   `json:"status"`
		DueDate   int64 `json:"due_date"`
		PayType   int   `json:"vip_pay_type"`
		ThemeType int   `json:"theme_type"`
		Label     struct {
			Path                  string `json:"path"`
			Text                  string `json:"text"`
			LabelTheme            string `json:"label_theme"`
			TextColor             string `json:"text_color"`
			BgStyle               int    `json:"bg_style"`
			BgColor               string `json:"bg_color"`
			BorderColor           string `json:"border_color"`
			UseImgLabel           bool   `json:"use_img_label"`
			ImgLabelUriHans       string `json:"img_label_uri_hans"`
			ImgLabelUriHant       string `json:"img_label_uri_hant"`
			ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
			ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
		} `json:"label"`
		AvatarSubscript    int    `json:"avatar_subscript"`
		NicknameColor      string `json:"nickname_color"`
		Role               int    `json:"role"`
		AvatarSubscriptUrl string `json:"avatar_subscript_url"`
		TVVIPStatus        int    `json:"tv_vip_status"`
		TVVIPPayType       int    `json:"tv_vip_pay_type"`
		TVDueDate          int64  `json:"tv_due_date"`
	} `json:"vip"`
	Pendant struct { // 头像框信息
		Pid               int    `json:"pid"`
		Name              string `json:"name"`
		Image             string `json:"image"`
		Expire            int    `json:"expire"`
		ImageEnhance      string `json:"image_enhance"`
		ImageEnhanceFrame string `json:"image_enhance_frame"`
	} `json:"pendant"`
	Nameplate struct { // 勋章信息
		Nid        int    `json:"nid"`
		Name       string `json:"name"`
		Image      string `json:"image"`
		ImageSmall string `json:"image_small"`
		Level      string `json:"level"`
		Condition  string `json:"condition"`
	} `json:"nameplate"`
	UserHonourInfo struct { // ？
		Mid    int64    `json:"mid"`
		Colour []string `json:"colour"`
		Tags   []string `json:"tags"`
	} `json:"user_honour_info"`
	IsFollowed bool     `json:"is_followed"` // 是否关注此用户 true 已关注 false 未关注，需要登陆，未登陆恒为false
	TopPhoto   string   `json:"top_photo"`   // 主页头像链接
	Theme      struct{} `json:"theme"`       // ？
	SysNotice  struct{} `json:"sys_notice"`  // 系统通知
	LiveRoom   struct { // 直播间信息
		RoomStatus    int    `json:"roomStatus"`
		LiveStatus    int    `json:"liveStatus"`
		URL           string `json:"url"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		RoomId        int    `json:"roomid"`
		RoundStatus   int    `json:"roundStatus"`
		BroadcastType int    `json:"broadcast_type"`
		WatchedShow   struct {
			Switch       bool   `json:"switch"`
			Num          int    `json:"num"`
			TextSmall    string `json:"text_small"`
			TextLarge    string `json:"text_large"`
			Icon         string `json:"icon"`
			IconLocation string `json:"icon_location"`
			IconWeb      string `json:"icon_web"`
		} `json:"watched_show"`
	} `json:"live_room"`
	Birthday string   `json:"birthday"` // 生日，如设置为隐私为空
	School   struct { // 学校
		Name string `json:"name"`
	} `json:"school"`
	Profession struct { // 专业资质信息
		Name       string `json:"name"`
		Department string `json:"department"`
		Title      string `json:"title"`
		IsShow     int    `json:"is_show"`
	} `json:"profession"`
	Tags   interface{} `json:"tags"` // 个人标签
	Series struct {
		UserUpgradeStatus int  `json:"user_upgrade_status"`
		ShowUpgradeWindow bool `json:"show_upgrade_window"`
	} `json:"series"`
	IsSeniorMember int         `json:"is_senior_member"` // 是否为硬核会员 0 否 1 是
	MCNInfo        interface{} `json:"mcn_info"`         // ？
	GaiaResType    int         `json:"gaia_res_type"`    // ？
	GaiaData       interface{} `json:"gaia_data"`        // ？
	IsRisk         bool        `json:"is_risk"`          // ？
	Elec           struct {    // 充电信息
		ShowInfo struct {
			Show    bool   `json:"show"`
			State   int    `json:"state"`
			Title   string `json:"title"`
			Icon    string `json:"icon"`
			JumpURL string `json:"jump_url"`
		} `json:"show_info"`
	} `json:"elec"`
	Contract struct { // 是否显示老粉计划
		IsDisplay       bool `json:"is_display"`
		IsFollowDisplay bool `json:"is_follow_display"`
	} `json:"contract"`
	CertificateShow bool `json:"certificate_show"` // ？
}

// GetUserCardResponse 用户名片信息
type GetUserCardResponse struct {
	Card struct {
		Mid         string        `json:"mid"`           // mid
		Name        string        `json:"name"`          // 昵称
		Approve     bool          `json:"approve"`       // ？
		Sex         string        `json:"sex"`           // 性别 男 女 保密
		Rank        string        `json:"rank"`          // 等级
		Face        string        `json:"face"`          // 用户头像链接
		FaceNft     int           `json:"face_nft"`      // 是否是nft头像 0 否 1 是
		FaceNftType int           `json:"face_nft_type"` // nft头像类别？
		DisplayRank string        `json:"DisplayRank"`   // ？
		RegTime     int64         `json:"regtime"`       // ？
		SpaceSta    int           `json:"spacesta"`      // ？
		Birthday    string        `json:"birthday"`      // 空
		Place       string        `json:"place"`         // 空
		Description string        `json:"description"`   // 空
		Article     int           `json:"article"`       // 0
		Attentions  []interface{} `json:"attentions"`    // 空
		Fans        int           `json:"fans"`          // 粉丝数
		Friend      int           `json:"friend"`        // 关注数
		Attention   int           `json:"attention"`     // 关注数
		Sign        string        `json:"sign"`          // 签名
		LevelInfo   struct {      // 等级信息
			CurrentLevel int `json:"current_level"`
			CurrentMin   int `json:"current_min"`
			CurrentExp   int `json:"current_exp"`
			NextExp      int `json:"next_exp"`
		} `json:"level_info"`
		Pendant struct { // 挂件
			Pid               int    `json:"pid"`
			Name              string `json:"name"`
			Expire            int    `json:"expire"`
			ImageEnhance      string `json:"image_enhance"`
			ImageEnhanceFrame string `json:"image_enhance_frame"`
		} `json:"pendant"`
		Nameplate struct { // 勋章
			Nid        int    `json:"nid"`
			Name       string `json:"name"`
			Image      string `json:"image"`
			ImageSmall string `json:"image_small"`
			Level      string `json:"level"`
			Condition  string `json:"condition"`
		} `json:"nameplate"`
		Official struct { // 认证信息
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"Official"`
		OfficialVerify struct { // 认证信息2
			Type int    `json:"type"`
			Desc string `json:"desc"`
		} `json:"official_verify"`
		VIP struct { // 会员信息
			Type       int   `json:"type"`
			Status     int   `json:"status"`
			DueDate    int64 `json:"due_date"`
			VipPayType int   `json:"vip_pay_type"`
			ThemeType  int   `json:"theme_type"`
			Label      struct {
				Path                  string `json:"path"`
				Text                  string `json:"text"`
				LabelTheme            string `json:"label_theme"`
				TextColor             string `json:"text_color"`
				BgStyle               int    `json:"bg_style"`
				BgColor               string `json:"bg_color"`
				BorderColor           string `json:"border_color"`
				UseImgLabel           bool   `json:"use_img_label"`
				ImgLabelUriHans       string `json:"img_label_uri_hans"`
				ImgLabelUriHant       string `json:"img_label_uri_hant"`
				ImgLabelUriHansStatic string `json:"img_label_uri_hans_static"`
				ImgLabelUriHantStatic string `json:"img_label_uri_hant_static"`
			} `json:"label"`
			AvatarSubscript    int    `json:"avatar_subscript"`
			NicknameColor      string `json:"nickname_color"`
			Role               int    `json:"role"`
			AvatarSubscriptUrl string `json:"avatar_subscript_url"`
			TVVIPStatus        int    `json:"tv_vip_status"`
			TVVIPPayType       int    `json:"tv_vip_pay_type"`
			TVDueDate          int64  `json:"tv_due_date"`
			VIPType            int    `json:"vipType"`
			VIPStatus          int    `json:"vipStatus"`
		} `json:"vip"`
		IsSeniorMember int `json:"is_senior_member"`
	} `json:"card"` // 卡片信息
	Space struct { // 主页图像
		SImg string `json:"s_img"`
		LImg string `json:"l_img"`
	} `json:"space"`
	Following    bool `json:"following"`     // 是否关注此用户 需登陆
	ArchiveCount int  `json:"archive_count"` // 用户稿件数
	ArticleCount int  `json:"article_count"` // ？
	Follower     int  `json:"follower"`      // 粉丝数
	LikeNum      int  `json:"like_num"`      // 点赞数
}

// GetMyInfoResponse 登陆用户个人详细信息
type GetMyInfoResponse struct {
	Mid            int64  `json:"mid"`             // mid
	Name           string `json:"name"`            // 昵称
	Sex            string `json:"sex"`             // 性别 男 女 保密
	Face           string `json:"face"`            // 头像图片url
	Sign           string `json:"sign"`            // 签名
	Rank           int    `json:"rank"`            // 10000
	Level          int    `json:"level"`           // 当前等级 0-6
	JoinTime       int    `json:"jointime"`        // 0 ?
	Moral          int    `json:"moral"`           // 节操 默认70
	Silence        int    `json:"silence"`         // 封禁状态 0 正常 1 被封
	EmailStatus    int    `json:"email_status"`    // 已验证邮箱 0 未验证 1 已验证
	TelStatus      int    `json:"tel_status"`      // 已验证手机号 0 未验证 1 已验证
	Identification int    `json:"identification"`  // 1 ？
	Birthday       int64  `json:"birthday"`        // 生日
	IsTourist      int    `json:"is_tourist"`      // 0 ？
	IsFakeAccount  int    `json:"is_fake_account"` // 0 ？
	PinPrompting   int    `json:"pin_prompting"`   // 0 ？
	IsDeleted      int    `json:"is_deleted"`      // 0 ？
	InRegAudit     int    `json:"in_reg_audit"`    // ？
	IsRipUser      bool   `json:"is_rip_user"`     // ？
	Profession     struct {
		ID              int    `json:"id"`
		Name            string `json:"name"`
		ShowName        string `json:"show_name"`
		IsShow          int    `json:"is_show"`
		CategoryOne     string `json:"category_one"`
		RealName        string `json:"realname"`
		Title           string `json:"title"`
		Department      string `json:"department"`
		CertificateNo   string `json:"certificate_no"`
		CertificateShow bool   `json:"certificate_show"`
	} `json:"profession"`
	FaceNft        int `json:"face_nft"`
	FaceNftNew     int `json:"face_nft_new"`
	IsSeniorMember int `json:"is_senior_member"`
	Honours        struct {
		Mid    int64 `json:"mid"`
		Colour struct {
			Dark   string `json:"dark"`
			Normal string `json:"normal"`
		} `json:"colour"`
		Tags interface{} `json:"tags"`
	} `json:"honours"`
	DigitalID   string `json:"digital_id"`
	DigitalType int    `json:"digital_type"`
	Attestation struct {
		Type       int `json:"type"`
		CommonInfo struct {
			Title       string `json:"title"`
			Prefix      string `json:"prefix"`
			PrefixTitle string `json:"prefix_title"`
		} `json:"common_info"`
		SpliceInfo struct {
			Title string `json:"title"`
		} `json:"splice_info"`
		Icon string `json:"icon"`
		Desc string `json:"desc"`
	} `json:"attestation"`
	ExpertInfo struct {
		Title string `json:"title"`
	} `json:"expert_info"`
	LevelExp struct { // 等级经验
		CurrentLevel int   `json:"current_level"`
		CurrentMin   int   `json:"current_min"`
		CurrentExp   int   `json:"current_exp"`
		NextExp      int   `json:"next_exp"`
		LevelUp      int64 `json:"level_up"`
	} `json:"level_exp"`
	Coins     float64 `json:"coins"`     // 硬币
	Following int     `json:"following"` // 粉丝数
	Follower  int     `json:"follower"`  // 粉丝数
}

// GetRelationStatResponse 用户关系状态
type GetRelationStatResponse struct {
	Mid       int64 `json:"mid"`
	Following int   `json:"following"` // 关注数
	Whisper   int   `json:"whisper"`   // 悄悄关注数 需要登陆
	Black     int   `json:"black"`     // 黑名单数 需要登陆
	Follower  int   `json:"follower"`  // 粉丝数
}

// GetUpStatResponse up主状态
type GetUpStatResponse struct {
	Archive struct {
		EnableVT int `json:"enable_vt"`
		View     int `json:"view"` // 视频播放量
		VT       int `json:"vt"`
	} `json:"archive"`
	Article struct {
		View int `json:"view"` // 专栏阅读量
	} `json:"article"`
	Likes int `json:"likes"` // 点赞量
}

// GetDocUploadCountResponse 相簿投稿数
type GetDocUploadCountResponse struct {
	AllCount   int `json:"all_count"`   // 相簿总数 以下3个之和
	DrawCount  int `json:"draw_count"`  // 发布绘画数
	PhotoCount int `json:"photo_count"` // 发布摄影数
	DailyCount int `json:"daily_count"` // 发布日常（图片动态）数
}

// RelationUserResponse 用户关系响应
type RelationUserResponse struct {
	List      []RelationUser `json:"list"` // 列表
	ReVersion interface{}    `json:"re_version"`
	Total     int            `json:"total"` // 总数
}

type RelationUser struct {
	Mid            int64    `json:"mid"`
	Attribute      int      `json:"attribute"`     // 0 未关注 1 已关注 2 已关注 6 已互粉 128 已拉黑
	Mtime          int      `json:"mtime"`         // 关注对方时间
	Tag            []int    `json:"tag"`           // 分组ID
	Special        int      `json:"special"`       // 特别关注标志 0 否 1 是
	ContractInfo   struct{} `json:"contract_info"` // unknown
	Uname          string   `json:"uname"`
	Face           string   `json:"face"`
	Sign           string   `json:"sign"`
	FaceNft        int      `json:"face_nft"`
	OfficialVerify struct {
		Type  int    `json:"type"` // 1 已认证 -1 无认证
		Desc  string `json:"desc"`
		Role  int    `json:"role"`
		Title string `json:"title"`
	} `json:"official_verify"`
	Vip struct {
		VipType       int    `json:"vipType"`
		VipDueDate    int64  `json:"vipDueDate"`
		DueRemark     string `json:"dueRemark"`
		AccessStatus  int    `json:"accessStatus"`
		VipStatus     int    `json:"vipStatus"`
		VipStatusWarn string `json:"vipStatusWarn"`
		ThemeType     int    `json:"themeType"`
		Label         struct {
			Path        string `json:"path"`
			Text        string `json:"text"`
			LabelTheme  string `json:"label_theme"`
			TextColor   string `json:"text_color"`
			BgStyle     int    `json:"bg_style"`
			BgColor     string `json:"bg_color"`
			BorderColor string `json:"border_color"`
		} `json:"label"`
		AvatarSubscript    int    `json:"avatar_subscript"`
		NicknameColor      string `json:"nickname_color"`
		AvatarSubscriptUrl string `json:"avatar_subscript_url"`
	} `json:"vip"`
	NftIcon   string `json:"nft_icon"`
	RecReason string `json:"rec_reason"`
	TrackId   string `json:"track_id"`
}

// BatchModifyRelationResponse 批量操作关系
type BatchModifyRelationResponse struct {
	FailedFids []string `json:"failed_fids"` // 操作失败的 mid 列表
}

type Attribute int

const (
	// UnFollowed 未关注
	UnFollowed Attribute = 0
	// Followed 已关注
	Followed Attribute = 2
	// FollowEachOther 已互粉
	FollowEachOther Attribute = 6
	// InBlacklist 已拉黑
	InBlacklist Attribute = 128
)

// Relation 关系
type Relation struct {
	Mid       int64     `json:"mid"`
	Attribute Attribute `json:"attribute"`
	MTime     int64     `json:"mtime"` // 关注对方时间
	Tag       []int     `json:"tag"`
	Special   int       `json:"special"` // 1 特别关注
}

// AccRelation 相互关系
type AccRelation struct {
	Relation   Relation `json:"relation"`
	BeRelation Relation `json:"be_relation"`
}

// RelationTag 分组标签
type RelationTag struct {
	TagId int    `json:"tagid"` // 0 默认分组 -10 特别关注
	Name  string `json:"name"`  // 分组名称
	Count int    `json:"count"` // 分组成员数
	Tip   string `json:"tip"`   // 提示信息
}

// CreateRelationTagResponse 创建分组
type CreateRelationTagResponse struct {
	TagId int `json:"tagid"`
}

// LogoutResponse 登出
type LogoutResponse struct {
	RedirectUrl string `json:"redirectUrl"`
}

// CookieInfo ...
type CookieInfo struct {
	Refresh   bool  `json:"refresh"`
	Timestamp int64 `json:"timestamp"`
}

// RefreshCookieResponse ...
type RefreshCookieResponse struct {
	Status       int    `json:"status"`
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token"`
}

// ExpReward 每日经验奖励状态
type ExpReward struct {
	Login        bool `json:"login"`         // 每日登陆 true 已完成 false 未完成 完成奖励5经验
	Watch        bool `json:"watch"`         // 每日观看 true 已完成 false 未完成 完成奖励5经验
	Coins        int  `json:"coins"`         // 每日投币所奖励的经验 上限50 注：该值更新存在延迟 大概延迟几秒钟
	Share        bool `json:"share"`         // 每日分享 true 已完成 false 未完成 完成奖励5经验
	Email        bool `json:"email"`         // 绑定邮箱 false 未完成 true 已完成 首次完成奖励20经验
	Tel          bool `json:"tel"`           // 绑定手机号 false 未完成 true 已完成 首次完成奖励100经验
	SafeQuestion bool `json:"safe_question"` // 设置密保问题 false 未完成 true 已完成 首次完成奖励30经验
	IdentifyCard bool `json:"identify_card"` // 实名认证 false 未完成 true 已完成 首次完成奖励50经验
}

// TripleVideoResponse ...
type TripleVideoResponse struct {
	Like     bool `json:"like"`     // 是否点赞成功
	Coin     bool `json:"coin"`     // 是否投币成功
	Fav      bool `json:"fav"`      // 是否收藏成功
	Multiply int  `json:"multiply"` // 投币数量
}

type Video struct {
	Aid         int         `json:"aid"`       // avid
	Videos      int         `json:"videos"`    // 分P总数
	Tid         int         `json:"tid"`       // 分区ID
	Tname       string      `json:"tname"`     // 子分区名称
	Copyright   int         `json:"copyright"` // 视频类型 1 原创 2 转载
	Pic         string      `json:"pic"`       // 封面url
	Title       string      `json:"title"`     // 标题
	Pubdate     int64       `json:"pubdate"`   // 发布时间 秒级时间戳
	Ctime       int64       `json:"ctime"`     // 投稿时间 秒级时间戳
	Desc        string      `json:"desc"`      // 简介
	DescV2      []*DescV2   `json:"desc_v2"`   // 新版视频简介
	State       int         `json:"state"`     // 状态
	Duration    int         `json:"duration"`  // 所有分P总时长 单位秒
	Rights      *Rights     `json:"rights"`    // 视频属性标志
	Owner       *Owner      `json:"owner"`     // Up主信息
	Stat        Stat        `json:"stat"`      // 视频状态数
	Dynamic     string      `json:"dynamic"`   // 视频同步发布的的动态的文字内容
	Cid         int64       `json:"cid"`       // 视频1P cid
	Dimension   *Dimension  `json:"dimension"` // 视频1P分辨率
	Pages       []*Page     `json:"pages"`     // 视频分P列表
	ShortLink   string      `json:"short_link_v2"`
	UpFromV2    int         `json:"up_from_v2"`
	FirstFrame  string      `json:"first_frame"`
	PubLocation string      `json:"pub_location"`
	Bvid        string      `json:"bvid"` // bvid
	SeasonType  int         `json:"season_type"`
	IsOgv       bool        `json:"is_ogv"`
	OgvInfo     interface{} `json:"ogv_info"`
	EnableVt    int         `json:"enable_vt"`
}

type DescV2 struct {
	RawText string `json:"raw_text"` // 简介内容 type=1时显示原文 type=2时显示'@'+raw_text+' '并链接至biz_id的主页
	Type    int    `json:"type"`     // 类型 1：普通，2：@他人
	BizId   int64  `json:"biz_id"`   // 被@用户的mid	=0，当type=1
}

type Rights struct {
	Bp            int `json:"bp"`             // 是否允许承包
	Elec          int `json:"elec"`           // 是否支持充电
	Download      int `json:"download"`       // 是否允许下载
	Movie         int `json:"movie"`          // 是否电影
	Pay           int `json:"pay"`            // 是否PGC付费
	Hd5           int `json:"hd5"`            // 是否有高码率
	NoReprint     int `json:"no_reprint"`     // 是否显示“禁止转载”标志
	Autoplay      int `json:"autoplay"`       // 是否自动播放
	UgcPay        int `json:"ugc_pay"`        // 是否UGC付费
	IsCooperation int `json:"is_cooperation"` // 是否为联合投稿
	UgcPayPreview int `json:"ugc_pay_preview"`
	NoBackground  int `json:"no_background"`
	ArcPay        int `json:"arc_pay"`
	PayFreeWatch  int `json:"pay_free_watch"`
}

type Owner struct {
	Mid  int64  `json:"mid"`  // 用户 mid
	Name string `json:"name"` // 用户名
	Face string `json:"face"` // 用户头像
}

type Stat struct {
	Aid      int64 `json:"aid"`      // avid
	View     int   `json:"view"`     // 播放数
	Danmaku  int   `json:"danmaku"`  // 弹幕数
	Reply    int   `json:"reply"`    // 评论数
	Favorite int   `json:"favorite"` // 收藏数
	Coin     int   `json:"coin"`     // 投币数
	Share    int   `json:"share"`    // 分享数
	NowRank  int   `json:"now_rank"` // 当前排名
	HisRank  int   `json:"his_rank"` // 历史最高排行
	Like     int   `json:"like"`     // 点赞数
	Dislike  int   `json:"dislike"`  // 点踩数 恒为0
	Vt       int   `json:"vt"`
	Vv       int   `json:"vv"`
}

type Dimension struct {
	Width  int `json:"width"`  // 宽度
	Height int `json:"height"` // 高度
	Rotate int `json:"rotate"` // 是否将宽高对换 0 正常 1 对换
}

type Page struct {
	Cid       int64      `json:"cid"`       // 分p cid
	Page      int        `json:"page"`      // 分p序号
	From      string     `json:"from"`      // 视频来源 vupload 普通上传 hunan 芒果TV qq 腾讯
	Part      string     `json:"part"`      // 分P标题
	Duration  int        `json:"duration"`  // 分P持续时间 单位秒
	Vid       string     `json:"vid"`       // 站外视频vid
	Weblink   string     `json:"weblink"`   // 站外视频跳转url
	Dimension *Dimension `json:"dimension"` // 当前分P分辨率
}

// GetPopularVideoListResponse ...
type GetPopularVideoListResponse struct {
	List   []*Video `json:"list"`
	NoMore bool     `json:"no_more"`
}

// GetLatestVideoResponse ...
type GetLatestVideoResponse struct {
	Archives []*Video `json:"archives"`
	Page     struct {
		Count int `json:"count"` // 总数
		Num   int `json:"num"`   // 当前页码
		Size  int `json:"size"`  // 每页项数
	} `json:"page"`
}
