package models

type ZssqCategory struct {
	CategoryName  string `form:"category_name" json:"category_name"`
	CategoryAlias string `form:"category_alias" json:"category_alias"`
	BookCount     int    `form:"book_count"   json:"book_count"`
	Gender        string `form:"gender" json:"gender"`
	Use           int    `form:"use" json:"use"`
}

type ZssqBookDesc struct {
	BookKey          string `form:"book_key" json:"book_key"`
	BookName         string `form:"book_name" json:"book_name"`
	Author           string `form:"author"   json:"author"`
	Pic              string `form:"pic" json:"pic"`
	Desc             string `form:"desc" json:"desc"`
	Serialize        int    `form:"serialize"   json:"serialize"`
	CategoryName     string `form:"category_name" json:"category_name"`
	ClassId          int64  `form:"class_id" json:"class_id"`
	Tags             string `form:"tags" json:"tags"`
	TextNum          int    `form:"text_num"   json:"text_num"`
	LastChapterTitle string `form:"last_chapter_title" json:"last_chapter_title"`
	LastChapterTime  string `form:"last_chapter_time" json:"last_chapter_time"`
	ChapterNum       int    `form:"chapter_num"   json:"chapter_num"`
	BookType         int    `form:"book_type"   json:"book_type"`
	IsClassic        int    `form:"is_classic"   json:"is_classic"`
	Use              int    `form:"use" json:"use"`
}

type ZssqChapter struct {
	ChapterName string `form:"chapter_name" json:"chapter_name"`
	ChapterLink string `form:"chapter_link" json:"chapter_link"`
	Sort        int    `form:"sort"   json:"sort"`
	Text        string `form:"text"   json:"text"`
	TextNum     int    `form:"text_num"   json:"text_num"`
}
