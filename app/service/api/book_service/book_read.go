package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetReadCountByBookId(bookId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McBookRead{}).Distinct("uid").Where("bid = ?", bookId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getSearchCountByBookName(bookName string) (count int64) {
	var err error
	err = global.DB.Model(models.McBookSearch{}).Distinct("uid").Where("search_name = ?", bookName).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetReadByUserId(bookId, userId int64) (read *models.McBookRead, err error) {
	err = global.DB.Model(models.McBookRead{}).Where("bid = ? and uid = ?", bookId, userId).First(&read).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getBrowseCountByBookId(bookId, userId int64) (count int64) {
	err := global.DB.Model(models.McBookBrowse{}).Where("bid = ? and uid = ?", bookId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func getBrowseByBookId(bookId, userId int64) (browse *models.McBookBrowse) {
	err := global.DB.Model(models.McBookBrowse{}).Where("bid = ? and uid = ?", bookId, userId).Last(&browse).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func updateBrowseByUserId(userId, bookId int64) (err error) {
	data := make(map[string]interface{})
	data["uptime"] = utils.GetUnix()
	err = global.DB.Model(models.McBookBrowse{}).Where("uid = ? and bid = ?", userId, bookId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func updateShelfByUserId(userId, bookId int64) (err error) {
	data := make(map[string]interface{})
	data["uptime"] = utils.GetUnix()
	err = global.DB.Model(models.McBookShelf{}).Where("uid = ? and bid = ?", userId, bookId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
