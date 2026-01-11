package book_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetBookBuyCountByUnique(uid, bid, cid int64) (count int64) {
	err := global.DB.Model(models.McBookBuy{}).Where("uid = ? and bid = ? and cid = ?", uid, bid, cid).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetBookBuyAutoByUnique(uid, bid int64) (auto int) {
	err := global.DB.Model(models.McBookBuy{}).Where("uid = ? and bid = ? ", uid, bid).Find(&auto).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateBookBuyAuto(bookId int64, auto int) (err error) {
	err = global.DB.Model(models.McBookBuy{}).Where("id = ?", bookId).Update("auto", auto).Error
	if err != nil {
		global.Sqllog.Error(err.Error())
		return
	}
	return
}
