package shelf_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetTodayShelfCountByUserId(uid int64) (count int64) {
	err := global.DB.Model(models.McBookShelf{}).Where("uid = ? and addtime >= ?", uid, utils.GetTomorrowUnix()).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetShelfCountByBookId(bid, uid int64) (count int64) {
	err := global.DB.Model(models.McBookShelf{}).Where("bid = ? and uid = ?", bid, uid).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateShitsByBookId(bookId int64) (err error) {
	err = global.DB.Model(models.McBook{}).Where("id", bookId).Update("shits", gorm.Expr("shits + ?", 1)).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteBookShelf(bookIds []int64, userId int64) (res bool, err error) {
	err = global.DB.Where("bid in ? and uid = ?", bookIds, userId).Delete(&models.McBookShelf{}).Error
	if err != nil {
		return
	}
	return true, nil
}

func getShelfSecondByUserId(uid, starWeekUnix int64) (seconds int64) {
	err := global.DB.Model(models.McBookTime{}).Select("coalesce(sum(second), 0)").Where("uid = ? and addtime >= ?", uid, starWeekUnix).Scan(&seconds).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteShelfByUserId(userId int64) (err error) {
	if userId <= 0 {
		return
	}
	err = global.DB.Where("uid = ?", userId).Delete(&models.McBookShelf{}).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
