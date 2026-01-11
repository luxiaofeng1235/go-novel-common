package user_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func DeleteCionByUserId(userId int64) (err error) {
	if userId <= 0 {
		return
	}
	err = global.DB.Where("uid = ?", userId).Delete(&models.McCionChange{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
