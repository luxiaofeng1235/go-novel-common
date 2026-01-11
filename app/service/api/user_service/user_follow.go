package user_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetFollowCountByUid(userId, byUserId int64) (count int64) {
	err := global.DB.Model(models.McUserFollow{}).Where("uid = ? and by_uid = ?", userId, byUserId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteFollowUid(userId, byUserId int64) (err error) {
	err = global.DB.Where("uid = ? and by_uid = ?", userId, byUserId).Delete(&models.McUserFollow{}).Error
	if err != nil {
		return
	}
	return
}

func GetFollowIdsByUid(userId int64) (userIds []int64) {
	err := global.DB.Model(models.McUserFollow{}).Debug().Where("uid = ?", userId).Pluck("by_uid", &userIds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetFansIdsByUid(userId int64) (userIds []int64) {
	err := global.DB.Model(models.McUserFollow{}).Where("by_uid = ?", userId).Pluck("uid", &userIds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteFollowByUserId(userId int64) (err error) {
	if userId <= 0 {
		return
	}
	err = global.DB.Where("uid = ? or by_uid = ?", userId, userId).Delete(&models.McUserFollow{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
