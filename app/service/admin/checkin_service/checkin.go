package checkin_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetRewardById(id int64) (reward *models.McCheckinReward, err error) {
	err = global.DB.Model(models.McCheckinReward{}).Where("id", id).First(&reward).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func RewardListSearch(req *models.CheckinRewardListReq) (list []*models.McCheckinReward, total int64, err error) {
	db := global.DB.Model(&models.McCheckinReward{}).Order("id desc")

	day := strings.TrimSpace(req.Day)
	if day != "" {
		db = db.Where("day = ?", day)
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

func CreateReward(req *models.CreateRewardReq) (InsertId int64, err error) {
	day := req.Day
	if day <= 0 {
		err = fmt.Errorf("%v", "日期不能为空")
		return
	}
	cion := req.Cion
	vip := req.Vip
	if cion <= 0 && vip <= 0 {
		err = fmt.Errorf("%v", "奖励不能为空")
		return
	}
	reward := models.McCheckinReward{
		Day:  day,
		Cion: cion,
		Vip:  vip,
	}

	if err = global.DB.Create(&reward).Error; err != nil {
		return 0, err
	}

	return reward.Id, nil
}

func UpdateReward(req *models.UpdateRewardReq) (res bool, err error) {
	id := req.RewardId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	var mapData = make(map[string]interface{})
	mapData["day"] = req.Day
	mapData["cion"] = req.Cion
	mapData["vip"] = req.Vip
	if err = global.DB.Model(models.McCheckinReward{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteReward(req *models.DeleteRewardReq) (res bool, err error) {
	ids := req.RewardIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.McCheckinReward{}).Error
		if err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return
	}
	return true, nil
}

func CheckinListSearch(req *models.CheckinListSearchReq) (list []*models.McCheckin, total int64, err error) {
	db := global.DB.Model(&models.McCheckin{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	day := strings.TrimSpace(req.Day)
	if day != "" {
		db = db.Where("day = ?", day)
	}

	isReissue := strings.TrimSpace(req.IsReissue)
	if isReissue != "" {
		db = db.Where("is_reissue = ?", isReissue)
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
