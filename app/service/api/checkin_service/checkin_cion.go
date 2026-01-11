package checkin_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetCheckinCionList() (checkIns []*models.McCheckinReward, err error) {
	err = global.DB.Model(models.McCheckinReward{}).Order("day asc").Find(&checkIns).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCheckinCionByDay(day int) (cion *models.McCheckinReward, err error) {
	err = global.DB.Model(models.McCheckinReward{}).Where("day = ?", day).First(&cion).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
