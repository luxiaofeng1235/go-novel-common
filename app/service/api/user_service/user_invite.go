package user_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetInviteCountByDeviceid(deviceid string) (count int64) {
	err := global.DB.Model(models.McUserInvite{}).Where("deviceid = ?", deviceid).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetInviteCountByUid(uid int64) (count int64) {
	err := global.DB.Model(models.McUserInvite{}).Where("uid = ?", uid).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
