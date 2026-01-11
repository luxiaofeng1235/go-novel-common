package class_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetClassBySex() (classList []*models.McBookClass, err error) {
	classList, err = GetClassList()
	if err != nil {
		err = fmt.Errorf("获取分类列表失败 err=%v", err.Error())
		return
	}
	if len(classList) > 0 {
		for _, val := range classList {
			text := val.ClassName
			if val.BookType == 1 {
				text = fmt.Sprintf("男生-%v", text)
			} else if val.BookType == 2 {
				text = fmt.Sprintf("女生-%v", text)
			}
			val.ClassName = text
		}
	}
	return
}

func BookListSearch(req *models.BookListReq) (list []*models.McBook, total int64, err error) {
	db := global.DB.Model(&models.McBook{}).Order("id desc")

	bookName := strings.TrimSpace(req.BookName)
	if bookName != "" {
		db = db.Where("book_name LIKE ?", "%"+bookName+"%")
	}

	bookType := strings.TrimSpace(req.BookType)
	if bookType != "" {
		db = db.Where("book_type = ?", bookType)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, book := range list {
		book.Pic = utils.GetAdminFileUrl(book.Pic)
	}
	return
}

func getBookPicById(id int64) (pic string) {
	var err error
	err = global.DB.Model(models.McBook{}).Select("pic").Where("id", id).First(&pic).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
