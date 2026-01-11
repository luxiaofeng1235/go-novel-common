package findbook_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/book_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

// 查找待处理的找书
func GetFindBookByUid(userId int64) (find *models.McFindbook) {
	var err error
	err = global.DB.Model(models.McFindbook{}).Where("status = 0 and uid = ?", userId).Last(&find).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func List(req *models.ApiFindbookListReq) (list []*models.McFindbook, total int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	db := global.DB.Model(&models.McFindbook{}).Order("id desc")

	db = db.Where("uid = ?", userId)

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}
	if pageSize == 0 {
		pageSize = 1
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return
}

func CreateFindbook(req *models.CreateFindBookReq) (bookId int64, err error) {
	bookName := strings.TrimSpace(req.BookName)
	userId := req.UserId
	if bookName == "" {
		err = fmt.Errorf("%v", "小说名称不能为空")
		return
	}

	author := strings.TrimSpace(req.Author)
	if author == "" {
		err = fmt.Errorf("%v", "小说作者不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	bookId = book_service.GetBookIdByBookNameAndAuthor(bookName, author)
	if bookId > 0 {
		return
	}

	find := GetFindBookByUid(userId)
	if find.Id > 0 {
		err = fmt.Errorf("%v 求助资源整理中 请耐心等候", find.BookName)
		return
	}

	sourceName := strings.TrimSpace(req.SourceName)
	findbook := models.McFindbook{
		BookName:   bookName,
		Author:     author,
		SourceName: sourceName,
		Uid:        userId,
		Status:     0,
		Addtime:    utils.GetUnix(),
		Uptime:     utils.GetUnix(),
	}
	if err = global.DB.Create(&findbook).Error; err != nil {
		err = fmt.Errorf("提交失败，稍后再试 err=%v", err.Error())
		return
	}
	return
}
