package findbook_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetFindBookById(id int64) (findbook *models.McFindbook, err error) {
	err = global.DB.Model(models.McFindbook{}).Where("id", id).First(&findbook).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func FindBookListSearch(req *models.FindBookListReq) (list []*models.McFindbook, total int64, err error) {
	db := global.DB.Model(&models.McFindbook{}).Order("id desc")

	bookName := strings.TrimSpace(req.BookName)
	if bookName != "" {
		db = db.Where("book_name like ?", "%"+bookName+"%")
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}
	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err
}

func UpdateFindBook(req *models.UpdateFindBookReq) (res bool, err error) {
	id := req.FindBookId
	bookName := strings.TrimSpace(req.BookName)
	if bookName == "" {
		err = fmt.Errorf("%v", "小说名称不能为空")
		return
	}
	author := strings.TrimSpace(req.Author)
	var mapData = make(map[string]interface{})
	mapData["book_name"] = bookName
	if author != "" {
		mapData["author"] = author
	}
	sourceName := strings.TrimSpace(req.SourceName)
	if sourceName != "" {
		mapData["source_name"] = sourceName
	}
	mapData["status"] = req.Status
	mapData["uptime"] = utils.GetUnix()
	mapData["book_times"] = req.BookTimes //求书次数设置

	if err = global.DB.Model(models.McFindbook{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DelFindBook(req *models.DelFindBookReq) (res bool, err error) {
	ids := req.FindBookIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McFindbook{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
