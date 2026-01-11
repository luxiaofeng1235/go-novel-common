package withdraw_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetLimitById(id int64) (limit *models.McWithdrawLimit, err error) {
	err = global.DB.Model(models.McWithdrawLimit{}).Where("id", id).First(&limit).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func LimitListSearch(req *models.WithdrawLimitListReq) (list []*models.McWithdrawLimit, total int64, err error) {
	db := global.DB.Model(&models.McWithdrawLimit{}).Order("id desc")

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
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

func CreateLimit(req *models.CreateLimitReq) (InsertId int64, err error) {
	if req.Money <= 0 {
		err = fmt.Errorf("%v", "可提现金额不能为空")
		return
	}

	limit := models.McWithdrawLimit{
		Money:   req.Money,
		Read:    req.Read,
		Status:  req.Status,
		Sort:    req.Sort,
		Addtime: utils.GetUnix(),
	}

	if err = global.DB.Create(&limit).Error; err != nil {
		return 0, err
	}

	return limit.Id, nil
}

func UpdateLimit(req *models.UpdateLimitReq) (res bool, err error) {
	id := req.LimitId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["money"] = req.Money
	mapData["read"] = req.Read
	mapData["status"] = req.Status
	mapData["sort"] = req.Sort
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McWithdrawLimit{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteLimit(req *models.DeleteLimitReq) (res bool, err error) {
	ids := req.LimitIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McWithdrawLimit{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}
