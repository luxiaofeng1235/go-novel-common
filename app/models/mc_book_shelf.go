package models

type McBookShelf struct {
	Id      int64 `gorm:"column:id" json:"id"`
	Uid     int64 `gorm:"column:uid" json:"uid"`         // 用户ID
	Bid     int64 `gorm:"column:bid" json:"bid"`         // 小说ID
	Top     int   `gorm:"column:top" json:"top"`         // 小说ID
	Addtime int64 `gorm:"column:addtime" json:"addtime"` // 添加时间
	Uptime  int64 `gorm:"column:uptime" json:"uptime"`   // 更新时间
}

func (*McBookShelf) TableName() string {
	return "mc_book_shelf"
}

type BookShelfListReq struct {
	Size   int   `form:"size" json:"size"`
	Page   int   `form:"page" json:"page"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type BookShelfListRes struct {
	Author             string `form:"author" json:"author"`
	Bid                int64  `form:"bid" json:"bid"`
	Uid                int64  `form:"uid" json:"uid"`
	Pic                string `form:"pic" json:"pic"`
	BookName           string `form:"book_name" json:"book_name"`
	NewsChapterId      int64  `form:"news_chapter_id" json:"news_chapter_id"`
	NewsChapterName    string `form:"news_chapter_name" json:"news_chapter_name"`
	ReadChapterId      int64  `form:"read_chapter_id" json:"read_chapter_id"`
	ReadChapterName    string `form:"read_chapter_name" json:"read_chapter_name"`
	ReadChapterTextNum int64  `form:"read_chapter_text_num" json:"read_chapter_text_num"`
	ChapterNum         int    `form:"chapter_num" json:"chapter_num"`
	Serialize          int    `form:"serialize" json:"serialize"`
	Top                int    `form:"top" json:"top"`
}

type BookShelfAddReq struct {
	BookId int64 `form:"bid" json:"bid"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type BookShelfDelReq struct {
	BookIds []int64 `form:"bids" json:"bids"`
	UserId  int64   `form:"user_id" json:"user_id"`
}

type BookShelfTopReq struct {
	BookId int64 `form:"bid" json:"bid"`
	Top    int   `form:"top" json:"top"`
	UserId int64 `form:"user_id" json:"user_id"`
}

type IsBookShelfReq struct {
	BookId int64 `form:"bid" json:"bid"`
	UserId int64 `form:"user_id" json:"user_id"`
}
