package withdraw_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetWithdrawLimitById(limitId int64) (limit *models.McWithdrawLimit, err error) {
	err = global.DB.Model(models.McWithdrawLimit{}).Where("id", limitId).First(&limit).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetAccountById(id int64) (account *models.McWithdrawAccount, err error) {
	err = global.DB.Model(models.McWithdrawAccount{}).Where("id", id).First(&account).Error
	return
}

func GetAccountByPayType(userId int64, payType int) (account *models.McWithdrawAccount, err error) {
	err = global.DB.Model(models.McWithdrawAccount{}).Where("uid = ? and pay_type = ?", userId, payType).First(&account).Error
	return
}

func GetAccountCountById(accountId, userId int64) (count int64) {
	var err error
	err = global.DB.Model(models.McWithdrawAccount{}).Where("id = ? and uid = ?", accountId, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}
