package withdraw_service

import (
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func WithdrawAccountListSearch(req *models.WithdrawAccountListReq) (list []*models.McWithdrawAccount, total int64, err error) {
	db := global.DB.Model(&models.McWithdrawAccount{}).Order("id desc")

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("uid = ?", userId)
	}

	cardName := req.CardName
	if cardName != "" {
		db = db.Where("card_name = ?", cardName)
	}

	cardNumber := req.CardNumber
	if cardNumber != "" {
		db = db.Where("card_number = ?", cardNumber)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		val.CardPic = utils.GetAdminFileUrl(val.CardPic)
	}
	return list, total, err
}
