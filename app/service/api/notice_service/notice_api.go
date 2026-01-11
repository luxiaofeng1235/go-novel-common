package notice_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetLastNotice() (notice *models.McNotice, err error) {
	err = global.DB.Model(&models.McNotice{}).Order("id desc").Where("status = 1").Last(&notice).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func MessageListSearch(req *models.MessageListReq, UserId int64) (list []*models.McNotify, total int64, pageNum, pageSize int, err error) {

	db := global.DB.Model(&models.McNotify{}).Order("id desc")

	db = db.Where("receive_uid = ? and notify_type not in ? ", UserId, []string{"praise", "comment"})
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

func ReplyList(req *models.ReplyListReq) (list []*models.McNotify, total int64, pageNum, pageSize int, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	db := global.DB.Model(&models.McNotify{}).Order("id desc")

	db = db.Where("receive_uid = ? and notify_type=?", userId, "comment")

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

	for _, val := range list {
		val.SendPic = utils.GetFileUrl(val.SendPic)
	}
	return
}

func PraiseList(req *models.PraiseListReq) (list []*models.McNotify, total int64, pageNum, pageSize int, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	db := global.DB.Model(&models.McNotify{}).Order("id desc")

	db = db.Where("receive_uid = ? and notify_type = ?", userId, utils.Praise)

	db = db.Where("addtime >= ?", utils.GetAgoDayUnix(180))

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
		val.SendPic = utils.GetFileUrl(val.SendPic)
	}
	return
}

func UpdateIsRead(req *models.UpdateIsReadReq) (err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	notifyId := req.NotifyId
	notifyType := strings.TrimSpace(req.MessageType)
	err = UpdateIsReadNotify(userId, notifyId, notifyType)
	return
}
