package withdraw_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func WithdrawListSearch(req *models.WithdrawListReq) (list []*models.McWithdrawApply, total int64, err error) {
	db := global.DB.Model(&models.McWithdrawApply{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
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

	for _, val := range list {
		val.CardPic = utils.GetAdminFileUrl(val.CardPic)
	}
	return list, total, err
}

func WithdrawCheck(req *models.WithdrawCheckReq) (err error) {
	mapData := make(map[string]interface{})
	mapData["status"] = req.Status
	if req.Status != 1 {
		mapData["reason"] = req.Reason
	}
	mapData["check_time"] = utils.GetUnix()
	if err = global.DB.Model(models.McWithdrawApply{}).Debug().Where("id", req.CheckId).Updates(&mapData).Error; err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return nil
}
