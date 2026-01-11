package withdraw_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetWithdrawAccountById(accountId int64) (account *models.McWithdrawAccount, err error) {
	err = global.DB.Model(models.McWithdrawAccount{}).Where("id", accountId).First(&account).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetCountByPayType(payType int, userId int64) (count int64) {
	err := global.DB.Model(models.McWithdrawAccount{}).Where("pay_type = ? and uid = ?", payType, userId).Count(&count).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func UpdateAccountByPayType(payType int, userId int64, data map[string]interface{}) (err error) {
	err = global.DB.Model(models.McWithdrawAccount{}).Where("pay_type = ? and uid = ?", payType, userId).Updates(data).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func DeleteAccountByAccountId(accountId int64, userId int64) (err error) {
	if accountId <= 0 || userId <= 0 {
		return
	}
	info, err := GetWithdrawAccountById(accountId)
	if err != nil {
		err = fmt.Errorf("%v", "账户不存在")
		return
	}
	err = global.DB.Where("id = ? and uid = ?", accountId, userId).Delete(&models.McWithdrawAccount{}).Error
	if err == nil {
		_ = utils.RemoveFile(info.CardPic)
	}
	return
}
