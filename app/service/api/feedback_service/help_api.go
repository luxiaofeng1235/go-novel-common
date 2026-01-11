package feedback_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"html"
)

func HelpListSearch(req *models.HelpListReq) (list []*models.McFeedbackHelp, total int64, pageNum, pageSize int, err error) {
	db := global.DB.Model(&models.McFeedbackHelp{}).Order("id desc")
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

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.Content = html.UnescapeString(val.Content)
	}
	return
}
