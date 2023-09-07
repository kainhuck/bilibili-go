package bilibili_go

import "github.com/kainhuck/bilibili-go/internal/utils"

/*
	视频分区：https://socialsisteryi.github.io/bilibili-API-collect/docs/video/video_zone.html
*/

type VideoZone struct {
	Name string
	Code string
	TID  int
	Desc string
	Url  string
}

type VideoZoneGroup struct {
	mainZone *VideoZone
	subZones map[int]*VideoZone
}

// GetMainZone 返回主分区
func (group *VideoZoneGroup) GetMainZone() *VideoZone {
	return group.mainZone
}

// GetSubZones 返回所有子分区
func (group *VideoZoneGroup) GetSubZones() map[int]*VideoZone {
	return group.subZones
}

// MainTid 返回分区的主tid
func (group *VideoZoneGroup) MainTid() int {
	return group.mainZone.TID
}

// RandomTid 返回一个随机分区tid
func (group *VideoZoneGroup) RandomTid() int {
	tids := make([]int, 0, len(group.subZones)+1)
	for tid := range group.subZones {
		tids = append(tids, tid)
	}
	tids = append(tids, group.MainTid())

	return utils.RandomChoice(tids)
}

// GetVideoZone 查询分区详细信息, 如果该分区不存在该tid则会返回主分区
func (group *VideoZoneGroup) GetVideoZone(tid int) *VideoZone {
	zone, ok := group.subZones[tid]
	if ok {
		return zone
	}

	return group.mainZone
}

// DougaGroup 动画分区
var DougaGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"动画(主分区)", "dougua", 1, "", "/v/douga"},
	subZones: map[int]*VideoZone{
		24:  {"MAD.AMV", "mad", 24, "具有一定创作度的动/静画二次创作视频", "/v/douga/mad"},
		25:  {"MMD·3D", "mmd", 24, "使用mmd（mikumikudance）和其他3d建模类软件制作的视频", "/v/douga/mmd"},
		47:  {"短片·手书·配音", "voice", 47, "追求个人特色和创意表达的自制动画短片、手书（绘）及acgn相关配音", "/v/douga/voice"},
		210: {"手办·模玩", "garage_kit", 210, "手办模玩的测评、改造或其他衍生内容", "/v/douga/garage_kit"},
		86:  {"特摄", "tokusatsu", 86, "特摄相关衍生视频", "/v/douga/tokusatsu"},
		253: {"动漫杂谈", "acgntalks", 253, "以谈话形式对ACGN文化圈进行的鉴赏、吐槽、评点、解说、推荐、科普等内容", "/v/douga/acgntalks"},
		27:  {"综合", "other", 27, "以动画及动画相关内容为素材，包括但不仅限于音频替换、恶搞改编、排行榜等内容", "/v/douga/other"},
	},
}

// AnimeGroup 番剧分区
var AnimeGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"番剧(主分区)", "anime", 13, "", "/anime"},
	subZones: map[int]*VideoZone{
		51:  {"资讯", "information", 51, "以动画/轻小说/漫画/杂志为主的资讯内容，PV/CM/特报/冒头/映像/预告", "/v/anime/information"},
		152: {"官方延伸", "official", 152, "以动画番剧及声优为主的EVENT/生放送/DRAMA/RADIO/LIVE/特典/冒头等", "/v/anime/official"},
		32:  {"完结动画", "finish", 32, "已完结TV/WEB动画及其独立系列，旧剧场版/OVA/SP/未放送", "/v/anime/finish"},
		33:  {"连载动画", "serial", 33, "连载中TV/WEB动画，新剧场版/OVA/SP/未放送/小剧场", "/v/anime/serial"},
	},
}

// GuochuangGroup 国创分区
var GuochuangGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"国创(主分区)", "guochuang", 167, "", "/guochuang"},
	subZones: map[int]*VideoZone{
		153: {"国产动画", "chinese", 153, "国产连载动画，国产完结动画", "/v/guochuang/chinese"},
		168: {"国产原创相关", "original", 168, "以国产动画、漫画、小说为素材的二次创作", "/v/guochuang/original"},
		169: {"布袋戏", "puppetry", 169, "布袋戏以及相关剪辑节目", "/v/guochuang/puppetry"},
		170: {"资讯", "information", 170, "原创国产动画、漫画的相关资讯、宣传节目等", "/v/guochuang/information"},
		195: {"动态漫·广播剧", "motioncomic", 195, "国产动态漫画、有声漫画、广播剧", "/v/guochuang/motioncomic"},
	},
}

// MusicGroup 音乐分区
var MusicGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"音乐(主分区)", "music", 3, "", "/v/music"},
	subZones: map[int]*VideoZone{
		28:  {"原创音乐", "original", 28, "原创歌曲及纯音乐，包括改编、重编曲及remix", "/v/music/original"},
		31:  {"翻唱", "cover", 31, "对曲目的人声再演绎视频", "/v/music/cover"},
		30:  {"VOCALOID·UTAU", "vocaloid", 30, "以vocaloid等歌声合成引擎为基础，运用各类音源进行的创作", "/v/music/vocaloid"},
		59:  {"演奏", "perform", 59, "乐器和非传统乐器器材的演奏作品", "/v/music/perform"},
		193: {"MV", "mv", 193, "为音乐作品配合拍摄或制作的音乐录影带（music video），以及自制拍摄、剪辑、翻拍mv", "/v/music/mv"},
		29:  {"音乐现场", "live", 29, "音乐表演的实况视频，包括官方/个人拍摄的综艺节目、音乐剧、音乐节、演唱会等", "/v/music/live"},
		130: {"音乐综合", "other", 130, "所有无法被收纳到其他音乐二级分区的音乐类视频", "/v/music/other"},
		243: {"乐评盘点", "commentary", 243, "音乐类新闻、盘点、点评、reaction、榜单、采访、幕后故事、唱片开箱等", "/v/music/commentary"},
		244: {"音乐教学", "tutorial", 244, "以音乐教学为目的的内容", "/v/music/tutorial"},
	},
}

// DanceGroup 舞蹈分区
var DanceGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"舞蹈(主分区)", "dance", 129, "", "/v/dance"},
	subZones: map[int]*VideoZone{
		20:  {"宅舞", "otaku", 20, "与acg相关的翻跳、原创舞蹈", "/v/dance/otaku"},
		154: {"舞蹈综合", "three_d", 154, "收录无法定义到其他舞蹈子分区的舞蹈视频", "/v/dance/three_d"},
		156: {"舞蹈教程", "demo", 156, "镜面慢速，动作分解，基础教程等具有教学意义的舞蹈视频", "/v/dance/demo"},
		198: {"街舞", "hiphop", 198, "收录街舞相关内容，包括赛事现场、舞室作品、个人翻跳、freestyle等", "/v/dance/hiphop"},
		199: {"明星舞蹈", "star", 199, "国内外明星发布的官方舞蹈及其翻跳内容", "/v/dance/star"},
		200: {"中国舞", "china", 200, "传承中国艺术文化的舞蹈内容，包括古典舞、民族民间舞、汉唐舞、古风舞等", "/v/dance/china"},
	},
}

// GameGroup 游戏分区
var GameGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"游戏(主分区)", "game", 4, "", "/v/game"},
	subZones: map[int]*VideoZone{
		17:  {"单机游戏", "stand_alone", 17, "以所有平台（pc、主机、移动端）的单机或联机游戏为主的视频内容，包括游戏预告、cg、实况解说及相关的评测、杂谈与视频剪辑等", "/v/game/stand_alone"},
		171: {"电子竞技", "esports", 171, "具有高对抗性的电子竞技游戏项目，其相关的赛事、实况、攻略、解说、短剧等视频", "/v/game/esports"},
		172: {"手机游戏", "mobile", 172, "以手机及平板设备为主要平台的游戏，其相关的实况、攻略、解说、短剧、演示等视频", "/v/game/mobile"},
		65:  {"网络游戏", "online", 65, "由网络运营商运营的多人在线游戏，以及电子竞技的相关游戏内容。包括赛事、攻略、实况、解说等相关视频", "/v/game/online"},
		173: {"桌游棋牌", "board", 173, "桌游、棋牌、卡牌对战等及其相关电子版游戏的实况、攻略、解说、演示等视频", "/v/game/board"},
		121: {"GMV", "gmv", 121, "由游戏素材制作的mv视频。以游戏内容或cg为主制作的，具有一定创作程度的mv类型的视频", "/v/game/gmv"},
		136: {"音游", "music", 136, "各个平台上，通过配合音乐与节奏而进行的音乐类游戏视频", "/v/game/music"},
		19:  {"Mugen", "mugen", 19, "以Mugen引擎为平台制作、或与Mugen相关的游戏视频", "/v/game/mugen"},
	},
}

// KnowledgeGroup 知识分区
var KnowledgeGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"知识(主分区)", "knowledge", 36, "", "/v/knowledge"},
	subZones: map[int]*VideoZone{
		201: {"科学科普", "science", 201, "回答你的十万个为什么", "/v/knowledge/science"},
		124: {"社科·法律·心理", "social_science", 124, "基于社会科学、法学、心理学展开或个人观点输出的知识视频", "/v/knowledge/social_science"},
		228: {"人文历史", "humanity_history", 228, "看看古今人物，聊聊历史过往，品品文学典籍", "/v/knowledge/humanity_history"},
		207: {"财经商业", "business", 207, "说金融市场，谈宏观经济，一起畅聊商业故事", "/v/knowledge/finance"},
		208: {"校园学习", "campus", 208, "老师很有趣，学生也有才，我们一起搞学习", "/v/knowledge/campus"},
		209: {"职业职场", "career", 209, "职业分享、升级指南，一起成为最有料的职场人", "/v/knowledge/career"},
		229: {"设计·创意", "design", 229, "天马行空，创意设计，都在这里", "/v/knowledge/design"},
		122: {"野生技术协会", "skill", 122, "技能党集合，是时候展示真正的技术了", "/v/knowledge/skill"},
	},
}

// TechGroup 科技分区
var TechGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"科技(主分区)", "tech", 188, "", "/v/tech"},
	subZones: map[int]*VideoZone{
		95:  {"数码", "digital", 95, "科技数码产品大全，一起来做发烧友", "/v/tech/digital"},
		230: {"软件应用", "application", 230, "超全软件应用指南", "/v/tech/application"},
		231: {"计算机技术", "computer_tech", 231, "研究分析、教学演示、经验分享......有关计算机技术的都在这里", "/v/tech/computer_tech"},
		232: {"科工机械", "industry", 232, "前方高能，机甲重工即将出没", "/v/tech/industry"},
		233: {"极客DIY", "diy", 233, "炫酷技能，极客文化，硬核技巧，准备好你的惊讶", "/v/tech/diy"},
	},
}

// SportsGroup 运动分区
var SportsGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"运动(主分区)", "sports", 234, "", "/v/sports"},
	subZones: map[int]*VideoZone{
		235: {"篮球", "basketball", 235, "与篮球相关的视频，包括但不限于篮球赛事、教学、评述、剪辑、剧情等相关内容", "/v/sports/basketball"},
		249: {"足球", "football", 249, "与足球相关的视频，包括但不限于足球赛事、教学、评述、剪辑、剧情等相关内容", "/v/sports/football"},
		164: {"健身", "aerobics", 164, "与健身相关的视频，包括但不限于瑜伽、crossfit、健美、力量举、普拉提、街健等相关内容", "/v/sports/aerobics"},
		236: {"竞技体育", "athletic", 236, "与竞技体育相关的视频，包括但不限于乒乓、羽毛球、排球、赛车等竞技项目的赛事、评述、剪辑、剧情等相关内容", "/v/sports/culture"},
		237: {"运动文化", "culture", 237, "与运动文化相关的视频，包络但不限于球鞋、球衣、球星卡等运动衍生品的分享、解读，体育产业的分析、科普等相关内容", "/v/sports/culture"},
		238: {"运动综合", "comprehensive", 238, "与运动综合相关的视频，包括但不限于钓鱼、骑行、滑板等日常运动分享、教学、Vlog等相关内容", "/v/sports/comprehensive"},
	},
}

// CarGroup 汽车分区
var CarGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"汽车(主分区)", "car", 223, "", "/v/car"},
	subZones: map[int]*VideoZone{
		245: {"赛车", "racing", 245, "f1等汽车运动相关", "/v/car/racing"},
		246: {"改装玩车", "modifiedvehicle", 246, "汽车文化及改装车相关内容，包括改装车、老车修复介绍、汽车聚会分享等内容", "/v/car/modifiedvehicle"},
		247: {"新能源车", "newenergyvehicle", 247, "新能源汽车相关内容，包括电动汽车、混合动力汽车等车型种类，包含不限于新车资讯、试驾体验、专业评测、技术解读、知识科普等内容", "/v/car/newenergyvehicle"},
		248: {"房车", "touringcar", 248, "房车及营地相关内容，包括不限于产品介绍、驾驶体验、房车生活和房车旅行等内容", "/v/car/touringcar"},
		240: {"摩托车", "motorcycle", 240, "骑士们集合啦", "/v/car/motorcycle"},
		227: {"购车攻略", "strategy", 227, "丰富详实的购车建议和新车体验", "/v/car/strategy"},
		176: {"汽车生活", "life", 176, "分享汽车及出行相关的生活体验类视频", "/v/car/life"},
	},
}

// LifeGroup 生活分区
var LifeGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"生活(主分区)", "life", 160, "", "/v/life"},
	subZones: map[int]*VideoZone{
		138: {"搞笑", "funny", 138, "各种沙雕有趣的搞笑剪辑，挑战，表演，配音等视频", "/v/life/funny"},
		250: {"出行", "travel", 250, "为达到观光游览、休闲娱乐为目的的远途旅行、中近途户外生活、本地探店", "/v/life/travel"},
		251: {"三农", "rurallife", 251, "分享美好农村生活", "/v/life/rurallife"},
		239: {"家居房产", "home", 239, "与买房、装修、居家生活相关的分享", "/v/life/home"},
		161: {"手工", "handmake", 161, "手工制品的制作过程或成品展示、教程、测评类视频", "/v/life/handmake"},
		162: {"绘画", "painting", 162, "绘画过程或绘画教程，以及绘画相关的所有视频", "/v/life/painting"},
		21:  {"日常", "daily", 21, "记录日常生活，分享生活故事", "/v/life/daily"},
	},
}

// FoodGroup 美食分区
var FoodGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"美食(主分区)", "food", 211, "", "/v/food"},
	subZones: map[int]*VideoZone{
		76:  {"美食制作", "make", 76, "学做人间美味，展示精湛厨艺", "/v/food/make"},
		212: {"美食侦探", "detective", 212, "寻找美味餐厅，发现街头美食", "/v/food/detective"},
		213: {"美食测评", "measurement", 213, "吃货世界，品尝世间美味", "/v/food/measurement"},
		214: {"田园美食", "rural", 214, "品味乡野美食，寻找山与海的味道", "/v/food/rural"},
		215: {"美食记录", "record", 215, "记录一日三餐，给生活添一点幸福感", "/v/food/record"},
	},
}

// AnimalGroup 动物圈分区
var AnimalGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"动物圈(主分区)", "animal", 217, "", "/v/animal"},
	subZones: map[int]*VideoZone{
		218: {"喵星人", "cat", 218, "喵喵喵喵喵", "/v/animal/cat"},
		219: {"汪星人", "dog", 219, "汪汪汪汪汪", "/v/animal/dog"},
		220: {"大熊猫", "panda", 220, "芝麻汤圆营业中", "/v/animal/panda"},
		221: {"野生动物", "wild_animal", 221, "内有“猛兽”出没", "/v/animal/wild_animal"},
		222: {"爬宠", "reptiles", 222, "鳞甲有灵", "/v/animal/reptiles"},
		75:  {"动物综合", "animal_composite", 75, "收录除上述子分区外，其余动物相关视频以及非动物主体或多个动物主体的动物相关延伸内容", "/v/animal/animal_composite"},
	},
}

// KichikuGroup 鬼畜分区
var KichikuGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"鬼畜(主分区)", "kichiku", 119, "", "/v/kichiku"},
	subZones: map[int]*VideoZone{
		22:  {"鬼畜调教", "guide", 22, "使用素材在音频、画面上做一定处理，达到与bgm一定的同步感", "/v/kichiku/guide"},
		26:  {"音MAD", "mad", 26, "使用素材音频进行一定的二次创作来达到还原原曲的非商业性质稿件", "/v/kichiku/mad"},
		126: {"人力VOCALOID", "manual_vocaloid", 126, "将人物或者角色的无伴奏素材进行人工调音，使其就像VOCALOID一样歌唱的技术", "/v/kichiku/manual_vocaloid"},
		216: {"鬼畜剧场", "theatre", 216, "使用素材进行人工剪辑编排的有剧情的作品", "/v/kichiku/theatre"},
		127: {"教程演示", "course", 127, "鬼畜相关的科普和教程演示", "/v/kichiku/course"},
	},
}

// FashionGroup 时尚分区
var FashionGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"时尚(主分区)", "fashion", 155, "", "/v/fashion"},
	subZones: map[int]*VideoZone{
		157: {"美妆护肤", "makeup", 157, "彩妆护肤、美甲美发、仿妆、医美相关内容分享或产品测评", "/v/fashion/makeup"},
		252: {"仿妆cos", "cos", 252, "对二次元、三次元人物角色进行模仿、还原、展示、演绎的内容", "/v/fashion/cos"},
		158: {"穿搭", "clothing", 158, "穿搭风格、穿搭技巧的展示分享，涵盖衣服、鞋靴、箱包配件、配饰（帽子、钟表、珠宝首饰）等", "/v/fashion/clothing"},
		159: {"时尚潮流", "catwalk", 159, "时尚街拍、时装周、时尚大片，时尚品牌、潮流等行业相关记录及知识科普", "/v/fashion/catwalk"},
	},
}

// InformationGroup 咨询分区
var InformationGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"资讯(主分区)", "information", 202, "", "/v/information"},
	subZones: map[int]*VideoZone{
		203: {"热点", "hotspot", 203, "全民关注的时政热门资讯", "/v/information/hotspot"},
		204: {"环球", "global", 204, "全球范围内发生的具有重大影响力的事件动态", "/v/information/global"},
		205: {"社会", "social", 205, "日常生活的社会事件、社会问题、社会风貌的报道", "/v/information/social"},
		206: {"综合", "multiple", 206, "除上述领域外其它垂直领域的综合资讯", "/v/information/multiple"},
	},
}

// EntGroup 娱乐分区
var EntGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"娱乐(主分区)", "ent", 5, "", "/v/ent"},
	subZones: map[int]*VideoZone{
		71:  {"综艺", "variety", 71, "所有综艺相关，全部一手掌握！", "/v/ent/variety"},
		241: {"娱乐杂谈", "talker", 241, "娱乐人物解读、娱乐热点点评、娱乐行业分析", "/v/ent/talker"},
		242: {"粉丝创作", "fans", 242, "粉丝向创作视频", "/v/ent/fans"},
		137: {"明星综合", "celebrity", 137, "娱乐圈动态、明星资讯相关", "/v/ent/celebrity"},
	},
}

// CinephileGroup 影视分区
var CinephileGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"影视(主分区)", "cinephile", 181, "", "/v/cinephile"},
	subZones: map[int]*VideoZone{
		182: {"影视杂谈", "cinecism", 182, "影视评论、解说、吐槽、科普等", "/v/cinephile/cinecism"},
		183: {"影视剪辑", "montage", 183, "对影视素材进行剪辑再创作的视频", "/v/cinephile/montage"},
		85:  {"小剧场", "shortfilm", 85, "有场景、有剧情的演绎类内容", "/v/cinephile/shortfilm"},
		184: {"预告·资讯", "trailer_info", 184, "影视类相关资讯，预告，花絮等视频", "/v/cinephile/trailer_info"},
	},
}

// DocumentaryGroup 纪录片分区
var DocumentaryGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"纪录片(主分区)", "documentary", 177, "", "/documentary"},
	subZones: map[int]*VideoZone{
		37:  {"人文·历史", "history", 37, "除宣传片、影视剪辑外的，人文艺术历史纪录剧集或电影、预告、花絮、二创、5分钟以上纪录短片", "/v/documentary/history"},
		178: {"科学·探索·自然", "science", 178, "除演讲、网课、教程外的，科学探索自然纪录剧集或电影、预告、花絮、二创、5分钟以上纪录短片", "/v/documentary/science"},
		179: {"军事", "military", 179, "除时政军事新闻外的，军事纪录剧集或电影、预告、花絮、二创、5分钟以上纪录短片", "/v/documentary/military"},
		180: {"社会·美食·旅行", "travel", 180, "除VLOG、风光摄影外的，社会美食旅行纪录剧集或电影、预告、花絮、二创、5分钟以上纪录短片", "/v/documentary/travel"},
	},
}

// MovieGroup 电影分区
var MovieGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"电影(主分区)", "movie", 23, "", "/movie"},
	subZones: map[int]*VideoZone{
		147: {"华语电影", "chinese", 147, "", "/v/movie/chinese"},
		145: {"欧美电影", "west", 145, "", "/v/movie/west"},
		146: {"日本电影", "japan", 146, "", "/v/movie/japan"},
		83:  {"其他国家", "movie", 83, "", "/v/movie/movie"},
	},
}

// TvGroup 电视剧分区
var TvGroup = &VideoZoneGroup{
	mainZone: &VideoZone{"电视剧(主分区)", "tv", 11, "", "/tv"},
	subZones: map[int]*VideoZone{
		185: {"国产剧", "mainland", 185, "", "/v/tv/mainland"},
		187: {"海外剧", "overseas", 187, "", "/v/tv/overseas"},
	},
}
