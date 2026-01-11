package models

type McBookRead struct {
	Id          int64  `gorm:"column:id" json:"id"`
	Uid         int64  `gorm:"column:uid" json:"uid"`                   // 用户ID
	Bid         int64  `gorm:"column:bid" json:"bid"`                   // 小说ID
	Cid         int64  `gorm:"column:cid" json:"cid"`                   // 章节ID
	ChapterName string `gorm:"column:chapter_name" json:"chapter_name"` // 章节名称
	TextNum     int64  `gorm:"column:text_num" json:"text_num"`         // 阅读到某章节多少字数
	Addtime     int64  `gorm:"column:addtime" json:"addtime"`           // 添加时间
	Uptime      int64  `gorm:"column:uptime" json:"uptime"`             // 更新时间
}

func (*McBookRead) TableName() string {
	return "mc_book_read"
}

type BookReadListReq struct {
	Size   int    `form:"size" json:"size"`
	Page   int    `form:"page" json:"page"`
	Day    string `form:"day" json:"day"`
	UserId int64  `form:"user_id" json:"user_id"`
}

type BookReadRes struct {
	Today     []*BookReadListRes `form:"today" json:"today"`
	Yesterday []*BookReadListRes `form:"yesterday" json:"yesterday"`
	Agoday    []*BookReadListRes `form:"agoday" json:"agoday"`
	Total     int64              `form:"total" json:"total"`
}

type BookReadListRes struct {
	Id              int64  `form:"id" json:"id"`
	Author          string `form:"author" json:"author"`
	BookName        string `form:"book_name" json:"book_name"`
	Pic             string `form:"pic" json:"pic"`
	Bid             int64  `form:"bid" json:"bid"`
	IsShelf         int64  `form:"is_shelf" json:"is_shelf"`
	TextNum         int64  `form:"text_num" json:"text_num"`
	NewsChapterId   int64  `form:"news_chapter_id" json:"news_chapter_id"`
	ReadChapterName string `form:"read_chapter_name" json:"read_chapter_name"`
	ReadChapterId   int64  `form:"read_chapter_id" json:"read_chapter_id"`
	ChapterNum      int    `form:"chapter_num" json:"chapter_num"`
	Serialize       int    `form:"serialize" json:"serialize"`
	Addtime         int64  `form:"addtime" json:"addtime"`
}

type ReadAddReq struct {
	BookId      int64  `form:"bid" json:"bid"`
	ChapterId   int64  `form:"cid" json:"cid"`
	ChapterName string `form:"chapter_name" json:"chapter_name"`
	TextNum     int64  `form:"text_num" json:"text_num"`
	Second      int64  `form:"second" json:"second"`
	UserId      int64  `form:"user_id" json:"user_id"`
}

type ReadInfoReq struct {
	BookId int64 `form:"bid" json:"bid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type ReadDelReq struct {
	BookIds []int64 `form:"bids" json:"bids"`
	UserId  int64   `form:"user_id" json:"user_id"`
}

type BrowseListReq struct {
	Size   int    `form:"size" json:"size"`
	Page   int    `form:"page" json:"page"`
	Day    string `form:"day" json:"day"`
	UserId int64  `form:"user_id" json:"user_id"`
}

type BrowseDelReq struct {
	BookIds []int64 `form:"bids" json:"bids"`
	UserId  int64   `form:"user_id" json:"user_id"`
}
