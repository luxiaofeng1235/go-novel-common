package chapter_service

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/global"
	"go-novel/utils"
)

func CreateChapter(bookName, author string, chapter *models.McBookChapter) (err error) {
	if chapter.ChapterName == "" {
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var chapterFile string
	chapterFile, err = GetChapterFile(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	var chapterAll []*models.McBookChapter
	chapterAll, err = GetChaptersByFile(chapterFile)

	var gq *gojsonq.JSONQ
	gq, err = GetBookJsonq(chapterFile)
	if err != nil {
		return
	}
	gq.Reset()
	count := gq.Where("chapter_name", "=", chapter.ChapterName).Count()
	if count > 0 {
		return
	}
	var lastIdChapter *models.McBookChapter
	lastIdChapter, err = GetLast(gq, "id")

	var lastSortChapter *models.McBookChapter
	lastSortChapter, err = GetLast(gq, "sort")

	var newId int64 = 1
	var newSort int = 1
	if lastIdChapter != nil && lastIdChapter.Id > 0 {
		newId = lastIdChapter.Id + 1
	}
	if lastIdChapter != nil && lastIdChapter.Sort > 0 {
		newSort = lastSortChapter.Sort + 1
	}

	newChapter := &models.McBookChapter{
		Id:          newId,
		Sort:        newSort,
		ChapterName: chapter.ChapterName,
		ChapterLink: chapter.ChapterLink,
		Vip:         0,
		Cion:        0,
		TextNum:     chapter.TextNum,
		Addtime:     utils.GetUnix(),
	}
	chapterAll = append(chapterAll, newChapter)
	newJsonData, err := json.MarshalIndent(chapterAll, "", "  ")
	if err != nil {
		err = fmt.Errorf("美化章节格式错误 %v", "")
		return
	}
	err = utils.WriteFile(chapterFile, string(newJsonData))
	return
}

func UpdateChapter(bookId int64, chapter *models.McBookChapter) (err error) {
	if chapter.ChapterName == "" {
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}
	var chapterFile string
	chapterFile, err = GetChapterFile(book.BookName, book.Author)
	if err != nil {
		return
	}
	var chapterAll []*models.McBookChapter
	chapterAll, err = GetChaptersByFile(chapterFile)

	var gq *gojsonq.JSONQ
	gq, err = GetBookJsonq(chapterFile)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	gq.Reset()

	count := gq.Where("chapter_name", "=", chapter.ChapterName).Count()
	if count <= 0 {
		return
	}

	newChapter := &models.McBookChapter{
		//Id:          newId,
		//Sort:        newSort,
		ChapterName: chapter.ChapterName,
		ChapterLink: chapter.ChapterLink,
		Vip:         0,
		Cion:        0,
		TextNum:     2000,
		Addtime:     utils.GetUnix(),
	}
	chapterAll = append(chapterAll, newChapter)
	newJsonData, err := json.MarshalIndent(chapterAll, "", "  ")
	if err != nil {
		err = fmt.Errorf("美化章节格式错误 %v", "")
		return
	}
	err = utils.WriteFile(chapterFile, string(newJsonData))
	return
}

func GetLast(gq *gojsonq.JSONQ, field string) (chapter *models.McBookChapter, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	gq.Reset()
	chapter = new(models.McBookChapter)
	last := gq.SortBy(field).Last()
	if last == nil {
		return
	}
	lastData, err := json.Marshal(last)
	if err != nil {
		return
	}
	err = json.Unmarshal(lastData, &chapter)
	if err != nil {
		return
	}
	return
}

func GetSortLast(bookName, author string) (sort int) {
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	var lastSortChapter = new(models.McBookChapter)
	lastSortChapter, err = GetLast(gq, "sort")
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	sort = lastSortChapter.Sort
	return
}

func GetFirst(gq *gojsonq.JSONQ, field string) (chapter *models.McBookChapter, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	gq.Reset()
	chapter = new(models.McBookChapter)
	first := gq.SortBy(field).First()
	if first == nil {
		return
	}
	lastData, err := json.Marshal(first)
	if err != nil {
		return
	}
	err = json.Unmarshal(lastData, &chapter)
	if err != nil {
		return
	}
	return
}

func GetSortFirst(bookName, author string) (sort int) {
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	var firstSortChapter = new(models.McBookChapter)
	firstSortChapter, err = GetFirst(gq, "sort")
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	sort = firstSortChapter.Sort
	return
}

func GetBookNewChapterId(bookName, author string) (newChapterId int64, newChapterName string) {
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	var lastSortChapter = new(models.McBookChapter)
	lastSortChapter, err = GetLast(gq, "sort")
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	newChapterId = lastSortChapter.Id
	newChapterName = lastSortChapter.ChapterName
	return
}

func GetChapterNameByChapterId(bookName, author string, chapterId int64) (chapterName string) {
	//章节表
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = GetJsonqByBookName(bookName, author)
	if err != nil {
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	gq.Reset()
	res := gq.Where("id", "=", chapterId).Get()
	if res == nil {
		return
	}
	info, err := json.Marshal(res)
	if err != nil {
		return
	}
	var chapter = new(models.McBookChapter)
	err = json.Unmarshal(info, &chapter)
	if err != nil {
		return
	}
	chapterName = chapter.ChapterName
	return
}

func GetBookNewChapterName(bookName, author string) (newChapterName string) {
	//章节表
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = GetJsonqByBookName(bookName, author)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	gq.Reset()
	res := gq.SortBy("sort", "desc").SortBy("id", "desc").Limit(1).First()
	info, err := json.Marshal(res)
	if err != nil {
		return
	}
	var chapter = new(models.McBookChapter)
	err = json.Unmarshal(info, &chapter)
	if err != nil {
		return
	}
	if chapter != nil {
		newChapterName = chapter.ChapterName
	}
	return
}

func GetChapterByChapterId(gq *gojsonq.JSONQ, chapterId int64) (chapter *models.McBookChapter, err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	gq.Reset()
	chapter = new(models.McBookChapter)
	res := gq.Where("id", "=", chapterId).First()
	if res == nil {
		return
	}
	info, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(info, &chapter)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	return
}

func GetChapterPrev(gq *gojsonq.JSONQ, sort int) (chapterId int64) {
	var err error
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var chapter = new(models.McBookChapter)
	gq.Reset()
	res := gq.SortBy("sort", "desc").Where("sort", "<", sort).Limit(1).First()
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	if res == nil {
		return
	}
	info, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(info, &chapter)
	if err != nil {
		return
	}
	chapterId = chapter.Id
	return
}

func GetChapterNext(gq *gojsonq.JSONQ, sort int) (chapterId int64) {
	var err error
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var chapter = new(models.McBookChapter)
	gq.Reset()
	res := gq.SortBy("sort", "asc").Where("sort", ">", sort).Limit(1).First()
	if res == nil {
		return
	}
	info, err := json.Marshal(res)
	if err != nil {
		return
	}
	err = json.Unmarshal(info, &chapter)
	if err != nil {
		return
	}
	chapterId = chapter.Id
	return
}
