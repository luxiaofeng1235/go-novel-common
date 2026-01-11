package models

type McBookChapter struct {
	Id          int64  `gorm:"id" json:"id"`
	Sort        int    `gorm:"sort" json:"sort"`                 // 排序ID
	ChapterLink string `gorm:"chapter_link" json:"chapter_link"` // 章节标题
	ChapterName string `gorm:"chapter_name" json:"chapter_name"` // 章节标题
	Vip         int    `gorm:"vip" json:"vip"`                   // VIP阅读，0否1是
	Cion        int64  `gorm:"cion" json:"cion"`                 // 章节需要金币
	IsFirst     int    `gorm:"is_first" json:"is_first"`         //是否为第一章
	IsLast      int    `gorm:"is_last" json:"is_last"`           //是否为最后一章
	TextNum     int    `gorm:"text_num" json:"text_num"`         // 章节字数
	Index       int    `gorm:"index" json:"index"`               // 索引
	IsLess      int    `gorm:"is_less" json:"is_less"`           //是否缺章
	Addtime     int64  `gorm:"addtime" json:"addtime"`           // 入库时间
}

func (*McBookChapter) TableName() string {
	return "mc_book_chapter"
}

type ChapterReq struct {
	BookId   int64  `form:"bid" json:"bid"`
	SourceId int64  `form:"source_id" json:"source_id"`
	Sort     string `form:"sort" json:"sort"` //asc-正序 desc-倒序
}

type ChapterBuyReq struct {
	BookId    int64  `form:"bid" json:"bid"` //小说ID
	ChapterId int64  `form:"zid" json:"zid"` //章节ID
	UserId    int64  `form:"user_id" json:"user_id"`
	Auto      int    `form:"auto" json:"auto"` //auto 1开启自动购买
	Ip        string `form:"ip" json:"ip"`
}

type ChapterReadReq struct {
	BookId    int64  `form:"bid" json:"bid"`
	ChapterId int64  `form:"cid" json:"cid"`
	UserId    int64  `form:"user_id" json:"user_id"`
	Ip        string `form:"ip" json:"ip"`
}

type ChapterListReq struct {
	BookId      int64  `form:"book_id" json:"book_id"`
	BookName    string `form:"book_name" json:"book_name"`
	Author      string `form:"author" json:"author"`
	ChapterId   string `form:"chapter_id" json:"chapter_id"`
	ChapterName string `form:"chapter_name" json:"chapter_name"`
	IsLess      string `form:"is_less" json:"is_less"`
	TextNumMin  int    `form:"text_num_min" json:"text_num_min"`
	TextNumMax  int    `form:"text_num_max" json:"text_num_max"`
	PageNum     int    `form:"pageNum" json:"pageNum"`
	PageSize    int    `form:"pageSize" json:"pageSize"`
}

type CreateChapterReq struct {
	BookId      int64  `form:"book_id"  json:"book_id"`
	BookName    string `form:"book_name" json:"book_name"`
	ChapterLink string `form:"chapter_link" json:"chapter_link"`
	Author      string `form:"author" json:"author"`
	ChapterName string `form:"chapter_name" json:"chapter_name"`
	TextNum     int    `form:"text_num" json:"text_num"`
	Sort        int    `form:"sort" json:"sort"`
	ChapterText string `form:"chapter_text" json:"chapter_text"`
}

type UpdateChapterReq struct {
	ChapterId int64 `form:"chapter_id"  json:"chapter_id"`
	Index     int   `form:"index"  json:"index"`
	CreateChapterReq
}

type DeleteChapterReq struct {
	BookId    int64 `form:"book_id" json:"book_id"`
	ChapterId int64 `form:"chapter_id" json:"chapter_id"`
}

type ChapterFeedBackAddReq struct {
	BookName     string `form:"book_name" json:"book_name"`
	Author       string `form:"author" json:"author"`
	ChapterName  string `form:"chapter_name" json:"chapter_name"`
	FeedbackType int    `form:"feedback_type" json:"feedback_type"`
	Bid          int64  `form:"bid" json:"bid"`
	Cid          int64  `form:"cid" json:"cid"`
	Ip           string `form:"ip" json:"ip"`
	UserId       int64  `form:"user_id" json:"user_id"`
}
