package models

type McBookSource struct {
	Id               int64  `gorm:"column:id" json:"id"`
	Bid              int64  `gorm:"column:bid" json:"bid"`                               // 小说ID
	Sid              int64  `gorm:"column:sid" json:"sid"`                               // 采集ID
	BookName         string `gorm:"column:book_name" json:"book_name"`                   // 小说名称
	Author           string `gorm:"column:author" json:"author"`                         // 小说作者
	SourceUrl        string `gorm:"column:source_url" json:"source_url"`                 // 小说采集链接地址
	LastChapterTitle string `gorm:"column:last_chapter_title" json:"last_chapter_title"` // 最新章节标题
	LastChapterTime  string `gorm:"column:last_chapter_time" json:"last_chapter_time"`   // 最新章节更新时间
	IsUpdate         int    `gorm:"column:is_update" json:"is_update"`                   // 是否继续更新 0-否 1-是
	Addtime          int64  `gorm:"column:addtime" json:"addtime"`                       // 添加时间
}

func (*McBookSource) TableName() string {
	return "mc_book_source"
}

type BookSourceReq struct {
	Bid int64 `form:"bid" json:"bid"`
}

type BookSourceRes struct {
	SourceId       int64  `form:"source_id" json:"source_id"`
	SourceUrl      string `form:"source_url" json:"source_url"`
	SourceName     string `form:"source_name" json:"source_name"`
	UpdateTime     string `form:"update_time" json:"update_time"`
	ChapterName    string `form:"chapter_name" json:"chapter_name"`
	ListSectionReg string `form:"list_section_reg" json:"list_section_reg"`
	ListUrlReg     string `form:"list_url_reg" json:"list_url_reg"`
	ChapterTextReg string `form:"chapter_text_reg" json:"chapter_text_reg"`
}

type SourceGetBookInfoReq struct {
	BookName string `form:"book_name" json:"book_name"`
	Author   string `form:"author" json:"author"`
}

type SourceGetBookInfoRes struct {
	BookName         string `form:"book_name" json:"book_name"`
	Pic              string `form:"pic" json:"pic"`
	Author           string `form:"author" json:"author"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize" json:"serialize"`
	ClassName        string `form:"class_name" json:"class_name"`
	Tags             string `form:"tags" json:"tags"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
}

type SourceUpdateBookInfoReq struct {
	BookName         string               `form:"book_name" json:"book_name"`
	Pic              string               `form:"pic" json:"pic"`
	Author           string               `form:"author" json:"author"`
	Desc             string               `form:"desc" json:"desc"`
	Serialize        int                  `form:"serialize" json:"serialize"`
	ClassName        string               `form:"class_name" json:"class_name"`
	Tags             string               `form:"tags" json:"tags"`
	LastChapterTitle string               `form:"last_chapter_title" json:"last_chapter_title"`
	Chapters         []*SourceChapterInfo `form:"chapters" json:"chapters"`
}

type SourceChapterInfo struct {
	ChapterTitle string `form:"chapter_title" json:"chapter_title"`
	ChapterLink  string `form:"chapter_link" json:"chapter_link"`
	TextNum      int    `form:"text_num" json:"text_num"`
	ChapterText  string `form:"chapter_text" json:"chapter_text"`
}

type NsqChapterInfoPush struct {
	BookId       int64  `form:"book_id" json:"book_id"`
	BookName     string `form:"book_name" json:"book_name"`
	Author       string `form:"author" json:"author"`
	ChapterTitle string `form:"chapter_title" json:"chapter_title"`
	ChapterLink  string `form:"chapter_link" json:"chapter_link"`
	TextNum      int    `form:"text_num" json:"text_num"`
	ChapterText  string `form:"chapter_text" json:"chapter_text"`
}
