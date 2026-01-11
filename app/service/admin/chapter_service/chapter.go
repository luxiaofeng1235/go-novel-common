package chapter_service

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/thedevsaddam/gojsonq/v2"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/app/service/common/chapter_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"io/ioutil"
	"strings"
)

func GetChapterById(bookName, author string, chapterId int64) (chapter *models.McBookChapter, err error) {
	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(bookName, author)
	if err != nil {
		return
	}
	var chapterAll []*models.McBookChapter
	chapterAll, err = chapter_service.GetChaptersByFile(chapterFile)
	if len(chapterAll) <= 0 {
		return
	}
	chapter = new(models.McBookChapter)
	for index, val := range chapterAll {
		if chapterId == val.Id {
			val.Index = index
			chapter = val
			break
		}
	}
	return
}

func ChapterListSearch(req *models.ChapterListReq) (list []*models.McBookChapter, total int, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "获取小说id失败")
		return
	}
	textNumMin := req.TextNumMin
	textNumMax := req.TextNumMax
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(req.BookName, req.Author)
	if err != nil {
		return
	}
	var gq *gojsonq.JSONQ
	gq, err = chapter_service.GetBookJsonq(chapterFile)
	if err != nil {
		return
	}
	gq.Reset()
	gq.SortBy("sort", "desc")
	// 当pageNum > 0 且 pageSize > 0 才分页
	pageNum := req.PageNum
	pageSize := req.PageSize

	chapterId := strings.TrimSpace(req.ChapterId)
	if chapterId != "" {
		gq.Where("id", "=", chapterId)
	}

	chapterName := strings.TrimSpace(req.ChapterName)
	if chapterName != "" {
		gq.Where("chapter_name", "contains", chapterName)
	}

	total = gq.Count()
	var res interface{}
	if pageNum > 0 && pageSize > 0 && req.IsLess == "" {
		res = gq.Offset((pageNum - 1) * pageSize).Limit(pageSize).Get()
	} else {
		res = gq.Get()
	}

	jsonData, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		err = fmt.Errorf("获取章节信息失败%v", err.Error())
		return
	}
	var tempList []*models.McBookChapter
	err = json.Unmarshal(jsonData, &tempList)
	if err != nil {
		err = fmt.Errorf("解析 JSON 数据错误 %v", err.Error())
		return
	}
	if req.IsLess == "1" {
		for _, val := range tempList {
			var txtFile string
			_, txtFile, err = book_service.GetChapterTxtFile(req.BookName, req.Author, val.ChapterName)
			if utils.CheckNotExist(txtFile) {
				list = append(list, val)
			}
		}
		total = len(list)
		if req.PageNum > 0 && req.PageSize > 0 {
			sliceStart, sliceEnd := utils.SlicePage(req.PageNum, req.PageSize, total)
			list = list[sliceStart:sliceEnd]
		}
	} else if req.IsLess == "0" {
		for _, val := range tempList {
			var txtFile string
			_, txtFile, err = book_service.GetChapterTxtFile(req.BookName, req.Author, val.ChapterName)
			if !utils.CheckNotExist(txtFile) {
				if textNumMin > 0 || textNumMax > 0 {
					textNum := book_service.GetTxtNum(txtFile)
					val.TextNum = textNum
					if textNum >= textNumMin && textNum <= textNumMax {
						list = append(list, val)
					}
				} else {
					list = append(list, val)
				}
			}
		}
		total = len(list)
		if req.PageNum > 0 && req.PageSize > 0 {
			sliceStart, sliceEnd := utils.SlicePage(req.PageNum, req.PageSize, total)
			list = list[sliceStart:sliceEnd]
		}
	} else {
		for _, val := range tempList {
			var txtFile string
			_, txtFile, err = book_service.GetChapterTxtFile(req.BookName, req.Author, val.ChapterName)
			if !utils.CheckNotExist(txtFile) {
				textNum := book_service.GetTxtNum(txtFile)
				val.TextNum = textNum
				if textNumMin > 0 || textNumMax > 0 {
					if textNum >= textNumMin && textNum <= textNumMax {
						list = append(list, val)
					}
				} else {
					list = append(list, val)
				}
			} else {
				list = append(list, val)
			}
		}
		if textNumMin > 0 || textNumMax > 0 {
			total = len(list)
		}
	}
	return
}

func CreateChapter(req *models.CreateChapterReq) (InsertId int64, err error) {
	chapterLink := strings.TrimSpace(req.ChapterLink)
	chapterName := strings.TrimSpace(req.ChapterName)
	if chapterName == "" {
		err = fmt.Errorf("%v", "章节名称不能为空")
		return
	}
	bookName := strings.TrimSpace(req.BookName)
	author := strings.TrimSpace(req.Author)
	if bookName == "" {
		err = fmt.Errorf("%v", "小说名称不能为空")
		return
	}
	if author == "" {
		err = fmt.Errorf("%v", "作者不能为空")
		return
	}

	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说id不正确")
		return
	}

	textNum := req.TextNum
	if textNum <= 0 {
		err = fmt.Errorf("%v", "章节字数不能为空")
		return
	}

	sort := req.Sort
	chapterText := req.ChapterText

	chapter := &models.McBookChapter{
		ChapterName: chapterName,
		ChapterLink: chapterLink,
		TextNum:     textNum,
		Sort:        sort,
		Addtime:     utils.GetUnix(),
	}

	err = chapter_service.CreateChapter(bookName, author, chapter)
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}

	_, chapterText, err = book_service.GetBookTxt(bookName, author, chapterName, chapterText)
	if err != nil {
		return
	}

	return chapter.Id, nil
}

func UpdateChapter(req *models.UpdateChapterReq) (res bool, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说id不正确")
		return
	}
	chapterId := req.ChapterId
	if chapterId <= 0 {
		err = fmt.Errorf("%v", "章节id不正确")
		return
	}
	chapterLink := strings.TrimSpace(req.ChapterLink)
	chapterName := strings.TrimSpace(req.ChapterName)
	bookName := strings.TrimSpace(req.BookName)
	author := strings.TrimSpace(req.Author)
	if bookName == "" {
		err = fmt.Errorf("%v", "小说名称不能为空")
		return
	}
	if author == "" {
		err = fmt.Errorf("%v", "作者不能为空")
		return
	}
	chapterText := req.ChapterText
	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(bookName, author)
	if err != nil {
		return
	}

	if chapterFile == "" {
		err = fmt.Errorf("%v", "小说章节json文件地址为空")
		return
	}
	var chapterAll []*models.McBookChapter
	chapterAll, err = chapter_service.GetChaptersByFile(chapterFile)
	if len(chapterAll) <= 0 {
		return
	}
	for _, val := range chapterAll {
		if chapterId == val.Id {
			val.Sort = req.Sort
			val.ChapterLink = chapterLink
			val.ChapterName = chapterName
			val.TextNum = req.TextNum
		}
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, err := json.MarshalIndent(chapterAll, "", "  ")
	if err != nil {
		err = fmt.Errorf("获取章节信息失败%v", err.Error())
		return
	}
	err = ioutil.WriteFile(chapterFile, jsonData, 0644)
	if err != nil {
		err = fmt.Errorf("更新章节错误 %v", err.Error())
		return
	}
	_, chapterText, err = book_service.GetBookTxt(bookName, author, chapterName, chapterText)
	if err != nil {
		return
	}
	return true, nil
}

func DeleteChapter(req *models.DeleteChapterReq) (res bool, err error) {
	bookId := req.BookId
	if bookId <= 0 {
		err = fmt.Errorf("%v", "小说id不正确")
		return
	}
	chapterId := req.ChapterId
	if chapterId <= 0 {
		err = fmt.Errorf("%v", "章节id不正确")
		return
	}
	book, err := book_service.GetBookById(bookId)
	if err != nil {
		return
	}
	var chapterFile string
	chapterFile, err = chapter_service.GetChapterFile(book.BookName, book.Author)
	if err != nil {
		return
	}

	if chapterFile == "" {
		err = fmt.Errorf("%v", "小说章节json文件地址为空")
		return
	}
	var chapterAll []*models.McBookChapter
	chapterAll, err = chapter_service.GetChaptersByFile(chapterFile)
	if len(chapterAll) <= 0 {
		return
	}

	var chapterNew []*models.McBookChapter
	var chapter *models.McBookChapter
	for index, val := range chapterAll {
		if chapterId == val.Id {
			chapterNew = append(chapterAll[:index], chapterAll[index+1:]...)
			chapter = &models.McBookChapter{
				Id:          val.Id,
				ChapterName: val.ChapterName,
				ChapterLink: val.ChapterLink,
			}
		}
	}

	var uploadBookTextPath string
	uploadBookTextPath, err = setting_service.GetValueByName(utils.UploadBookTextPath)
	if err != nil {
		err = fmt.Errorf("获取小说内容目录失败 uploadBookTextPath=%v", uploadBookTextPath)
		return
	}
	txtDir := fmt.Sprintf("%v%v", uploadBookTextPath, utils.GetBookMd5(book.BookName, book.Author))
	filePath := fmt.Sprintf("%v/%v.txt", txtDir, utils.GetChapterMd5(chapter.ChapterName))
	_ = utils.RemoveFile(filePath)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonData, err := json.MarshalIndent(chapterNew, "", "  ")
	if err != nil {
		err = fmt.Errorf("获取章节信息失败%v", err.Error())
		return
	}
	err = ioutil.WriteFile(chapterFile, jsonData, 0644)
	if err != nil {
		err = fmt.Errorf("更新章节错误 %v", err.Error())
		return
	}
	return true, nil
}

func GetSortLast(bookName, author string) (sort int) {
	var gq *gojsonq.JSONQ
	var err error
	gq, _, err = chapter_service.GetJsonqByBookName(bookName, author)
	if err != nil {
		return
	}
	var lastSortChapter *models.McBookChapter
	lastSortChapter, err = chapter_service.GetLast(gq, "sort")
	if err != nil {
		global.Jsonq.Errorf("%v", err.Error())
		return
	}
	sort = lastSortChapter.Sort
	return
}
