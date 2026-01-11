package vip_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetVipById(cardId int64) (vip *models.McVipCard, err error) {
	err = global.DB.Model(models.McVipCard{}).Where("id", cardId).First(&vip).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetUserVipById(id int64) (vip int64) {
	err := global.DB.Model(models.McUser{}).Select("vip").Where("id = ?", id).First(&vip).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetVipBookCountById() (count int64) {
	//Where("is_pay = 2") 暂时屏蔽掉
	err := global.DB.Model(models.McBook{}).Where("status=1").Debug().Count(&count).Error

	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getUserById(id int64) (user *models.McUser, err error) {
	err = global.DB.Model(models.McUser{}).Where("id", id).Find(&user).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
