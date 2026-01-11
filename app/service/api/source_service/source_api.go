package source_service

import (
	"fmt"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/api/collect_service"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/source_service"
	"go-novel/global"
	"go-novel/utils"
)

func GetBookUrlByBookId(sourceId, bookId int64) (sourceUrl string) {
	var err error
	err = global.DB.Model(models.McBookSource{}).Select("source_url").Where("sid = ? and bid = ?", sourceId, bookId).Scan(&sourceUrl).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
func GetBookSourceList(bookId int64) (sourceList []*models.BookSourceRes, err error) {
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说id不能为空")
		return
	}

	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}
	//章节表
	var gq *gojsonq.JSONQ
	gq, _, err = chapter_service.GetJsonqByBookName(book.BookName, book.Author)
	if err != nil {
		return
	}
	var chapterLast *models.McBookChapter
	chapterLast, err = chapter_service.GetLast(gq, "sort")
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}

	defaultSource := &models.BookSourceRes{
		SourceId:    0,
		SourceUrl:   book.SourceUrl,
		SourceName:  utils.DefaultSource,
		UpdateTime:  utils.UnixToDatetime(chapterLast.Addtime),
		ChapterName: chapterLast.ChapterName,
	}
	sourceList = append(sourceList, defaultSource)

	var bookSources []*models.McBookSource
	err = global.DB.Model(models.McBookSource{}).Where("bid = ?", bookId).Find(&bookSources).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if len(bookSources) <= 0 {
		return
	}

	for _, val := range bookSources {
		var collect *models.McCollect
		collect, err = collect_service.GetCollectById(val.Sid)
		if err != nil {
			continue
		}
		source := &models.BookSourceRes{
			SourceId:       val.Sid,
			SourceName:     collect.Title,
			ListSectionReg: collect.ChapterMode,
			ListUrlReg:     collect.ListUrlReg,
			ChapterTextReg: collect.ChapterTextReg,
			SourceUrl:      val.SourceUrl,
			ChapterName:    val.LastChapterTitle,
			UpdateTime:     val.LastChapterTime,
		}
		sourceList = append(sourceList, source)
		//go UpdateLastChapter(val, collect)
	}

	return
}

func UpdateLastChapter(bookSource *models.McBookSource, collect *models.McCollect) {
	var err error
	var updateTime, chapterTitle string
	_, updateTime, chapterTitle, err = source_service.GetSourceLastChapter(bookSource.SourceUrl, collect)
	if err != nil {
		return
	}
	data := make(map[string]interface{})
	data["last_chapter_title"] = chapterTitle
	data["last_chapter_time"] = updateTime
	data["uptime"] = utils.GetUnix()
	err = global.DB.Model(models.McBookSource{}).Where("id = ?", bookSource.Id).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
