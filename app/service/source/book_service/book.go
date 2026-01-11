package book_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func createBook(req *models.SourceUpdateBookInfoReq) (err error) {
	bookName := strings.TrimSpace(req.BookName)
	pic := strings.TrimSpace(req.Pic)
	author := strings.TrimSpace(req.Author)
	desc := strings.TrimSpace(req.Desc)
	className := strings.TrimSpace(req.ClassName)
	tags := strings.TrimSpace(req.Tags)
	lastChapterTitle := strings.TrimSpace(req.LastChapterTitle)
	serialize := req.Serialize

	if pic == "" {
		err = fmt.Errorf("%v", "小说图片不能为空")
		return
	}
	if author == "" {
		err = fmt.Errorf("%v", "小说作者不能为空")
		return
	}
	if desc == "" {
		err = fmt.Errorf("%v", "小说简介不能为空")
		return
	}
	if desc == "" {
		err = fmt.Errorf("%v", "小说简介不能为空")
		return
	}
	if serialize <= 0 {
		serialize = 1
	}
	if className == "" {
		err = fmt.Errorf("%v", "小说分类名称不能为空")
		return
	}
	if lastChapterTitle == "" {
		err = fmt.Errorf("%v", "最新章节不能为空")
		return
	}
	classId, bookType := getClassTypeByName(className)

	book := models.McBook{
		BookName:         bookName,
		Pic:              pic,
		Author:           author,
		Desc:             desc,
		Serialize:        serialize,
		ClassName:        className,
		Tags:             tags,
		Status:           1,
		Cid:              classId,
		BookType:         bookType,
		LastChapterTitle: lastChapterTitle,
		Addtime:          utils.GetUnix(),
	}
	if err = global.DB.Create(&book).Error; err != nil {
		return
	}
	return
}
