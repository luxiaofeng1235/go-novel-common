package book_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

//func UpdateBookInfo(req *models.SourceUpdateBookInfoReq) (err error) {
//	bookName := strings.TrimSpace(req.BookName)
//	if bookName == "" {
//		err = fmt.Errorf("%v", "小说书名不能为空")
//		return
//	}
//	pic := strings.TrimSpace(req.Pic)
//	author := strings.TrimSpace(req.Author)
//	desc := strings.TrimSpace(req.Desc)
//	className := strings.TrimSpace(req.ClassName)
//	tags := strings.TrimSpace(req.Tags)
//	lastChapterTitle := strings.TrimSpace(req.LastChapterTitle)
//	serialize := req.Serialize
//	count := book_service.GetBookCountByName(bookName)
//	if count <= 0 {
//		err = createBook(req)
//	} else {
//		var book *models.McBook
//		book, err = book_service.GetBookByBookName(bookName)
//		if err != nil {
//			return
//		}
//		if book.Id <= 0 {
//			err = createBook(req)
//			return
//		}
//		classId, bookType := getClassTypeByName(className)
//		data := make(map[string]interface{})
//		if bookName != "" {
//			data["book_name"] = bookName
//		}
//		if pic != "" {
//			data["pic"] = pic
//		}
//		if author != "" {
//			data["author"] = author
//		}
//		if desc != "" {
//			data["desc"] = desc
//		}
//		if serialize > 0 {
//			data["serialize"] = serialize
//		}
//		if className != "" {
//			data["class_name"] = className
//		}
//		if tags != "" {
//			data["tags"] = tags
//		}
//		if classId > 0 {
//			data["cid"] = classId
//		}
//		if bookType > 0 {
//			data["book_type"] = bookType
//		}
//		if bookType > 0 {
//			data["book_type"] = bookType
//		}
//		if lastChapterTitle != "" {
//			data["chapter_title"] = lastChapterTitle
//		}
//
//		err = global.DB.Model(models.McBook{}).Where("book_name = ?", bookName).Updates(data).Error
//		if err != nil {
//			global.Sqllog.Errorf("%v", err.Error())
//			return
//		}
//	}
//	return
//}

func GetBookInfo(req *models.SourceGetBookInfoReq) (res *models.SourceGetBookInfoRes, err error) {
	bookName := strings.TrimSpace(req.BookName)
	if bookName == "" {
		err = fmt.Errorf("%v", "小说书名不能为空")
		return
	}
	author := strings.TrimSpace(req.Author)
	if author == "" {
		err = fmt.Errorf("%v", "小说作者不能为空")
		return
	}
	var book *models.McBook
	book, err = book_service.GetBookByBookName(bookName, author)
	res = &models.SourceGetBookInfoRes{
		BookName:         book.BookName,
		Pic:              book.Pic,
		Author:           book.Author,
		Desc:             book.Desc,
		Serialize:        book.Serialize,
		ClassName:        book.ClassName,
		Tags:             book.Tags,
		LastChapterTitle: book.LastChapterTitle,
	}
	return
}

func UpdateBookChapter(req *models.SourceUpdateBookInfoReq) (err error) {
	bookName := strings.TrimSpace(req.BookName)
	if bookName == "" {
		err = fmt.Errorf("%v", "小说书名不能为空")
		return
	}
	author := strings.TrimSpace(req.Author)
	if author == "" {
		err = fmt.Errorf("%v", "小说作者不能为空")
		return
	}
	count := book_service.GetBookCountByName(bookName, author)
	if count <= 0 {
		err = createBook(req)
		if err != nil {
			return
		}
	}
	pic := strings.TrimSpace(req.Pic)
	desc := strings.TrimSpace(req.Desc)
	className := strings.TrimSpace(req.ClassName)
	tags := strings.TrimSpace(req.Tags)
	lastChapterTitle := strings.TrimSpace(req.LastChapterTitle)
	serialize := req.Serialize
	chapters := req.Chapters

	var book *models.McBook
	book, err = book_service.GetBookByBookName(bookName, author)
	bookId := book.Id
	if bookId <= 0 {
		err = fmt.Errorf("%v", "数据异常 该小说不存在 请稍后再试")
		return
	}

	classId, bookType := getClassTypeByName(className)
	data := make(map[string]interface{})
	if bookName != "" {
		data["book_name"] = bookName
	}
	if pic != "" {
		data["pic"] = pic
	}
	if author != "" {
		data["author"] = author
	}
	if desc != "" {
		data["desc"] = desc
	}
	if serialize > 0 {
		data["serialize"] = serialize
	}
	if className != "" {
		data["class_name"] = className
	}
	if tags != "" {
		data["tags"] = tags
	}
	if classId > 0 {
		data["cid"] = classId
	}
	if bookType > 0 {
		data["book_type"] = bookType
	}
	if bookType > 0 {
		data["book_type"] = bookType
	}
	if lastChapterTitle != "" {
		data["last_chapter_title"] = lastChapterTitle
	}

	err = global.DB.Model(models.McBook{}).Where("book_name = ?", bookName).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}

	for _, chapter := range chapters {
		msg := &models.NsqChapterInfoPush{
			BookId:       bookId,
			BookName:     bookName,
			Author:       author,
			ChapterTitle: chapter.ChapterTitle,
			ChapterLink:  chapter.ChapterLink,
			TextNum:      chapter.TextNum,
			ChapterText:  chapter.ChapterText,
		}
		var jsonData []byte
		jsonData, err = json.Marshal(msg)
		if err != nil {
			err = fmt.Errorf("转换json数据失败: %v", err.Error())
			return
		}
		// 发送当前批次的数据到 NSQ
		err = global.NsqPro.Publish(utils.SourceUpdateLastChapter, jsonData)
		if err != nil {
			err = fmt.Errorf("队列发送数据失败 %v", err.Error())
			return
		}
	}
	return
}
