package models

type McBook struct {
	Id                 int64   `gorm:"column:id" json:"id"`
	BookType           int     `gorm:"column:book_type" json:"book_type"`                       // 阅读类型
	BookName           string  `gorm:"column:book_name" json:"book_name"`                       // 小说书名
	Pic                string  `gorm:"column:pic" json:"pic"`                                   // 竖版封面
	Cid                int64   `gorm:"column:cid" json:"cid"`                                   // 分类ID
	IsRec              int     `gorm:"column:is_rec" json:"is_rec"`                             // 1-推荐，0-未推
	IsHot              int     `gorm:"column:is_hot" json:"is_hot"`                             // 是否热门
	IsChoice           int     `gorm:"column:is_choice" json:"is_choice"`                       // 精选 1-是，0-否
	IsClassic          int     `gorm:"column:is_classic" json:"is_classic"`                     // 经典 1-是，0-否
	IsNew              int     `gorm:"column:is_new" json:"is_new"`                             // 新书 1-是，0-否
	IsTeen             int     `gorm:"column:is_teen" json:"is_teen"`                           // 是否为青少年 1-是 0-否
	Serialize          int     `gorm:"column:serialize" json:"serialize"`                       // 0-连载 1-完结 3-太监
	Author             string  `gorm:"column:author" json:"author"`                             // 小说作者
	Tags               string  `gorm:"column:tags" json:"tags"`                                 // 采集书籍原标签
	Tid                int64   `gorm:"column:tid" json:"tid"`                                   // 标签ID
	TagName            string  `gorm:"column:tag_name" json:"tag_name"`                         // 标签名称
	CategoryName       string  `gorm:"column:category_name" json:"category_name"`               // 原分类名称
	Desc               string  `gorm:"column:desc" json:"desc"`                                 // 小说简介
	TextNum            int     `gorm:"column:text_num" json:"text_num"`                         // 总字数
	Hits               int64   `gorm:"column:hits" json:"hits"`                                 // 浏览数量
	HitsMonth          int64   `gorm:"column:hits_month" json:"hits_month"`                     // 月点击
	HitsWeek           int64   `gorm:"column:hits_week" json:"hits_week"`                       // 周点击
	HitsDay            int64   `gorm:"column:hits_day" json:"hits_day"`                         // 日点击
	Shits              int64   `gorm:"column:shits" json:"shits"`                               // 收藏人气
	IsPay              int     `gorm:"column:is_pay" json:"is_pay"`                             // 是否收费 0-免费 1-金币，2-VIP
	ChapterNum         int     `gorm:"column:chapter_num" json:"chapter_num"`                   // 章节总数
	Score              float64 `gorm:"column:score" json:"score"`                               // 总得分
	UpdateChapterId    int64   `gorm:"column:update_chapter_id" json:"update_chapter_id"`       // 最近一次更新章节id
	UpdateChapterTitle string  `gorm:"column:update_chapter_title" json:"update_chapter_title"` // 最近一次更新章节标题
	UpdateChapterTime  int64   `gorm:"column:update_chapter_time" json:"update_chapter_time"`   // 最近一次更新章节时间
	LastChapterTitle   string  `gorm:"column:last_chapter_title" json:"last_chapter_title"`     // 目标网站最新章节标题
	LastChapterTime    int64   `gorm:"column:last_chapter_time" json:"last_chapter_time"`       // 目标网站最新章节更新时间
	SourceUrl          string  `gorm:"column:source_url" json:"source_url"`                     // 采集来源标识
	SourceId           int64   `gorm:"column:source_id" json:"source_id"`                       // 采集资源ID
	ReadCount          int64   `gorm:"column:read_count" json:"read_count"`                     // 在读人数
	SearchCount        int64   `gorm:"column:search_count" json:"search_count"`                 // 搜索人数
	SearchNum          int64   `gorm:"column:search_num" json:"search_num"`                     // 搜索次数
	Status             int     `gorm:"column:status" json:"status"`                             // 分类状态,0-不展示,1-展示
	ClassName          string  `gorm:"column:class_name" json:"class_name"`                     // 小说分类名称 冗余
	Addtime            int64   `gorm:"column:addtime" json:"addtime"`                           // 入库时间
	Uptime             int64   `gorm:"column:uptime" json:"uptime"`                             // 更新时间
}

func (*McBook) TableName() string {
	return "mc_book"
}

// 搜索热门书籍的返回
type HotSearchRankRes struct {
	SearchNum int64 `form:"search_num" json:"search_num"`
	McBook
}

type SimpleBookRes struct {
	Id       int64  `form:"id" json:"id"`
	BookName string `form:"book_name" json:"book_name"` // 标题
	Author   string `form:"author" json:"author"`
	Pic      string `form:"pic" json:"pic"` // 竖版封面
}

type ApiCateBookReq struct {
	Bid         int64  `form:"bid" json:"bid"`                    //书籍ID
	Cid         int    `form:"cid" json:"cid"`                    //分类ID
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type ApiBookListReq struct {
	Cid         int    `form:"cid" json:"cid"`                     //分类ID
	Serialize   int    `form:"serialize" json:"serialize"`         //小说完本状态 1-连载 2-完结
	TextNumType int64  `form:"text_num_type" json:"text_num_type"` //1-50w以下 2-50:100w字 3-100万字以上
	Sort        string `form:"sort" json:"sort"`                   //排序 hot-热门 hits-人气 score-评分 new-最新
	IsPay       int    `form:"is_pay" json:"is_pay"`               //是否收费 0-免费 1-金币，2-VIP
	IsHot       int    `form:"is_hot" json:"is_hot"`               //是否热门  0：非热门 1：热门
	IsNew       int    `form:"is_new" json:"is_new"`               //是否为最新 0：非祖新 1：最新
	BookName    string `form:"book_name" json:"book_name"`
	Tid         int64  `form:"tid" json:"tid"`
	Page        int    `form:"page" json:"page"`
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type BookInfoReq struct {
	BookId int64 `form:"bid" json:"bid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type BookInfoRes struct {
	BookId          int64                   `form:"bid" json:"bid"`
	BookName        string                  `form:"book_name" json:"book_name"`
	Pic             string                  `form:"pic" json:"pic"`
	Serialize       int                     `form:"serialize" json:"serialize"`
	Author          string                  `form:"author" json:"author"`
	Desc            string                  `form:"desc" json:"desc"`
	TextNum         int                     `form:"text_num" json:"text_num"`
	ReadChapterId   int64                   `form:"read_chapter_id" json:"read_chapter_id"`
	ReadChapterName string                  `form:"read_chapter_name" json:"read_chapter_name"`
	IsPay           int                     `form:"is_pay" json:"is_pay"`
	Hits            int64                   `form:"hits" json:"hits"`
	Score           float64                 `form:"score" json:"score"`
	IsShelf         int                     `form:"is_shelf" json:"is_shelf"`
	CommentCount    int64                   `form:"comment_count" json:"comment_count"`
	ReadCount       int64                   `form:"read_count" json:"read_count"`
	Addtime         int64                   `form:"addtime" json:"addtime"`
	NewChapterName  string                  `form:"new_chapter_name" json:"new_chapter_name"`
	CommentList     []*CommentListRes       `form:"comment_list" json:"comment_list"`
	Scores          []*BookInfoHighScoreRes `form:"scores" json:"scores"`
}

type BookHighScoreReq struct {
	Bid         int64  `form:"bid" json:"bid"`
	ClassId     int64  `form:"cid" json:"cid"`
	TagId       int64  `form:"tid" json:"tid"`
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`
	Mark        string `form:"mark"  json:"mark"` //渠道号                  //所属IP，解析客户端IP用
}

type BookInfoHighScoreRes struct {
	Id       int64  `form:"id" json:"id"`
	BookName string `form:"book_name" json:"book_name"`
	Pic      string `form:"pic" json:"pic"`
}

type GetRankRes struct {
	Id   int64 `form:"id" json:"id"`
	Cion int   `form:"cion" json:"cion"`
	Rank int   `form:"rank" json:"rank"`
}

type ReadInfoRes struct {
	ChapterId     int64  `form:"chapter_id" json:"chapter_id"`
	ChapterName   string `form:"chapter_name" json:"chapter_name"`
	Bid           int64  `form:"bid" json:"bid"`
	IsShelf       int    `form:"is_shelf" json:"is_shelf"`
	IsAd          int    `form:"is_ad" json:"is_ad"`
	Vip           int    `form:"vip" json:"vip"`
	IsFirst       int    `form:"is_first" json:"is_first"`
	IsLast        int    `form:"is_last" json:"is_last"`
	PrevChapterId int64  `form:"prev_chapter_id" json:"prev_chapter_id"`
	NextChapterId int64  `form:"next_chapter_id" json:"next_chapter_id"`
	TextNum       int    `form:"text_num" json:"text_num"`
	Addtime       int64  `form:"addtime" json:"addtime"`
	Text          string `form:"text" json:"text"`
	AudioName     string `form:"audio_name" json:"audio_name"` //增加多媒体音频的播放数据返回
}

type BookScoreReq struct {
	BookId int64 `form:"bid" json:"bid"`
	Score  int   `form:"score" json:"score"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type VipBookStoreReq struct {
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type VipChoicesReq struct {
	Size        int    `form:"size" json:"size"`
	IsChoice    int    `form:"is_choice" json:"is_choice"`
	Page        int    `form:"page" json:"page"`
	IsHot       int    `form:"is_hot" json:"is_hot"`
	IsNew       int    `form:"is_new" json:"is_new"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type ApiRankListReq struct {
	BookType    int    `form:"book_type" json:"book_type"`
	ColumnType  int    `form:"column_type" json:"column_type"`
	Cid         int64  `form:"cid" json:"cid"`
	IsRec       int    `form:"is_rec" json:"is_rec"`
	IsHot       int    `form:"is_hot" json:"is_hot"`
	IsNew       int    `form:"is_new" json:"is_new"`
	IsClassic   int    `form:"is_classic" json:"is_classic"`
	Sort        string `form:"sort" json:"sort"`
	Tid         int64  `form:"tid" json:"tid"`
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type SectionForYouRecReq struct {
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type SectionHighScoreReq struct {
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type SectionEndReq struct {
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type SectionNewReq struct {
	Size        int    `form:"size" json:"size"`
	UserId      int64  `form:"user_id" json:"user_id"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type TagsRank struct {
	TagId    int64  `form:"tag_id" json:"tag_id"`
	TagName  string `form:"tag_name" json:"tag_name"`
	TagCount int64  `form:"tag_count" json:"tag_count"`
}

type TeenZoneListReq struct {
	TeenType    int    `form:"teen_type" json:"teen_type"` //1-名著 2-传记 3-外国文学
	Size        int    `form:"size" json:"size"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type GetTagsReq struct {
	ColumnType int   `form:"column_type" json:"column_type"`
	Size       int   `form:"size" json:"size"`
	UserId     int64 `form:"user_id" json:"user_id"`
}

type GetNewBookRecRes struct {
	TagId    int64            `form:"tag_id" json:"tag_id"`
	TagName  string           `form:"tag_name" json:"tag_name"`
	BookList []*SimpleBookRes `form:"book_list" json:"book_list"`
}

type GetNewBookListReq struct {
	Tid         int64  `form:"tid" json:"tid"`
	Page        int    `form:"page" json:"page"`
	Size        int    `form:"size" json:"size"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析
	Mark        string `form:"mark"  json:"mark"`
}

type ExchangeVipReq struct {
	VipCardId int64 `form:"vid" json:"vid"`
	UserId    int64 `form:"uid" json:"uid"`
}

type TodayUpdateBooksReq struct {
	Page        int    `form:"page" json:"page"`
	Size        int    `form:"size" json:"size"`
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type HotBookCountReq struct {
	DeviceType  string `form:"device_type"  json:"device_type"`   //设备类型  ios|android
	PackageName string `form:"package_name"  json:"package_name"` //包名
	Ip          string `form:"ip" json:"ip"`                      //所属IP，解析客户端IP用
	Mark        string `form:"mark"  json:"mark"`                 //渠道号
}

type BookListReq struct {
	BookType      string `form:"book_type" json:"book_type"`
	BookId        string `form:"book_id" json:"book_id"`
	BookName      string `form:"book_name" json:"book_name"`
	Author        string `form:"author" json:"author"`
	SourceUrl     string `form:"source_url" json:"source_url"`
	Serialize     string `form:"serialize" json:"serialize"`
	IsRec         string `form:"is_rec" json:"is_rec"`
	IsHot         string `form:"is_hot" json:"is_hot"`
	IsChoice      string `form:"is_choice" json:"is_choice"`
	IsClassic     string `form:"is_classic" json:"is_classic"`
	IsNew         string `form:"is_new" json:"is_new"`
	IsLess        string `form:"is_less" json:"is_less"`
	IsBanquan     string `form:"is_banquan" json:"is_banquan"`
	Status        string `form:"status" json:"status"`
	Cid           string `form:"cid" json:"cid"`
	Tid           string `form:"tid" json:"tid"`
	SourceDay     int    `form:"source_day" json:"source_day"`
	SourceNoday   int    `form:"source_noday" json:"source_noday"`
	RecentlyDay   int    `form:"recently_day" json:"recently_day"`
	EecentlyNoday int    `form:"recently_noday" json:"recently_noday"`
	BeginTime     string `form:"beginTime" json:"beginTime"`
	EndTime       string `form:"endTime" json:"endTime"`
	PageNum       int    `form:"pageNum" json:"pageNum"`
	PageSize      int    `form:"pageSize" json:"pageSize"`
}

// 定义审核通过的书籍列表
type BookListPassReq struct {
	BookName string `form:"book_name" json:"book_name"`
	Cid      string `form:"cid" json:"cid"`
	PageNum  int    `form:"pageNum" json:"pageNum"`
	PageSize int    `form:"pageSize" json:"pageSize"`
}

type BookListAdminRes struct {
	Id                 int64  `form:"id" json:"id"`
	BookName           string `form:"book_name" json:"book_name"`
	Author             string `form:"author" json:"author"`
	UpdateChapterTitle string `gorm:"column:update_chapter_title" json:"update_chapter_title"`
	SourceUrl          string `form:"source_url" json:"source_url"`
}

type CreateBookReq struct {
	BookType    int     `form:"book_type"  json:"book_type"`
	BookName    string  `form:"book_name" json:"book_name"`
	Pic         string  `form:"pic" json:"pic"`
	Cid         int64   `form:"cid" json:"cid"`
	IsRec       int     `form:"is_rec" json:"is_rec"`
	IsHot       int     `form:"is_hot" json:"is_hot"`
	IsChoice    int     `form:"is_choice" json:"is_choice"`
	IsClassic   int     `form:"is_classic" json:"is_classic"`
	IsNew       int     `form:"is_new" json:"is_new"`
	IsTeen      int     `form:"is_teen" json:"is_teen"`
	Serialize   int     `form:"serialize" json:"serialize"`
	Author      string  `form:"author" json:"author"`
	Tags        string  `form:"tags" json:"tags"`
	Desc        string  `form:"desc" json:"desc"`
	TextNum     int     `form:"text_num" json:"text_num"`
	Hits        int64   `form:"hits" json:"hits"`
	HitsMonth   int64   `form:"hits_month" json:"hits_month"`
	HitsWeek    int64   `form:"hits_week" json:"hits_week"`
	HitsDay     int64   `form:"hits_day" json:"hits_day"`
	Shits       int64   `form:"shits" json:"shits"`
	IsPay       int     `form:"is_pay" json:"is_pay"`
	ChapterNum  int     `form:"chapter_num" json:"chapter_num"`
	Score       float64 `form:"score" json:"score"`
	SourceUrl   string  `form:"source_url" json:"source_url"`
	SourceId    int64   `form:"source_id" json:"source_id"`
	ReadCount   int64   `form:"read_count" json:"read_count"`
	SearchCount int64   `form:"search_count" json:"search_count"`
	Status      int     `form:"status" json:"status"`
	IsIndex     int64   `form:"is_index" json:"is_index"` //判断是否为首页进入
}

type UpdateBookReq struct {
	BookId int64 `form:"id"  json:"id"`
	CreateBookReq
}

type DeleteBookReq struct {
	BookId int64 `json:"id" form:"id"`
}
