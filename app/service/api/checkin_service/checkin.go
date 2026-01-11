package checkin_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
)

func getUserVipTimeById(userId int64) (viptime int64) {
	var err error
	err = global.DB.Model(models.McUser{}).Select("viptime").Where("id", userId).First(&viptime).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getChekinHistory(userId int64, year, month int) (historys []*models.CheckinHistoryRes, err error) {
	tabName := new(models.McCheckin).TableName()
	sql := fmt.Sprintf("SELECT DAY(FROM_UNIXTIME(addtime)) AS day, cion, vip, is_reissue FROM %v WHERE uid = %v AND YEAR(FROM_UNIXTIME(addtime)) = %v AND MONTH(FROM_UNIXTIME(addtime)) = %v", tabName, userId, year, month)
	err = global.DB.Raw(sql).Scan(&historys).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCheckinUnixLastDay(userId, todayUnix int64) (dayUnix int64) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Order("id desc").Select("addtime").Where("is_reissue = 0 and uid = ? and addtime < ?", userId, todayUnix).First(&dayUnix).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCheckinLast(userId, todayUnix int64) (checkin *models.McCheckin) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Order("id desc").Where("is_reissue = 0 and uid = ? and addtime < ?", userId, todayUnix).First(&checkin).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCheckinLastDay(userId int64) (day int) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Order("id desc,day desc").Select("day").Where("uid = ?", userId).First(&day).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCountByDayTime(userId int64, startTime, endTime int64) (count int64) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Where("uid = ? and addtime >= ? and addtime <= ?", userId, startTime, endTime).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCountByDay(userId int64, day int) (count int64) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Where("uid = ? and day = ?", userId, day).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCountByAddtime(userId, addtime int64) (count int64) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Where("is_reissue = 0 and uid = ? and addtime >= ?", userId, addtime).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetTodayCountByUid(userId, todayUnix int64) (count int64) {
	var err error
	err = global.DB.Model(models.McCheckin{}).Where("is_reissue = 0 and uid = ? and addtime >= ?", userId, todayUnix).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
