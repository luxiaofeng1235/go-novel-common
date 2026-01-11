package feedback_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetReadCountByBookId(userId, agoTime int64) (count int64) {
	err := global.DB.Model(models.McFeedback{}).Where("uid = ? and addtime >= ?", userId, agoTime).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func FeedBackListSearch(req *models.FeedBackListReq) (list []*models.McFeedback, total int64, pageNum, pageSize int, err error) {

	db := global.DB.Model(&models.McFeedback{}).Order("id desc")
	userId := req.UserId
	if userId > 0 {
		db = db.Where("uid = ?", userId)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum = req.Page
	pageSize = req.Size

	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 || pageSize > 300 {
		pageSize = 15
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return
}
