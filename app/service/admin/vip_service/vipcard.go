package vip_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetCardById(id int64) (vipCard *models.McVipCard, err error) {
	err = global.DB.Model(models.McVipCard{}).Where("id", id).First(&vipCard).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func CheckCardNameUnique(tagName string, id int64) bool {
	var count int64
	model := global.DB.Model(models.McVipCard{}).Where("card_name = ?", tagName)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

func CardListSearch(req *models.CardListReq) (list []*models.McVipCard, total int64, err error) {
	db := global.DB.Model(&models.McVipCard{}).Order("id desc")

	cardName := strings.TrimSpace(req.CardName)
	if cardName != "" {
		db = db.Where("card_name = ?", cardName)
	}

	isRmb := strings.TrimSpace(req.IsRmb)
	if isRmb != "" {
		db = db.Where("is_rmb = ?", isRmb)
	}

	isCion := strings.TrimSpace(req.IsCion)
	if isCion != "" {
		db = db.Where("is_cion = ?", isCion)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

func CreateCard(req *models.CreateCardReq) (InsertId int64, err error) {
	cardName := strings.TrimSpace(req.CardName)
	if cardName == "" {
		err = fmt.Errorf("%v", "会员卡名称不能为空")
		return
	}

	if !CheckCardNameUnique(cardName, 0) {
		err = fmt.Errorf("%v", "会员卡名称已经存在")
		return
	}
	price := req.Price
	disRate := req.DisRate
	disPrice := req.DisPrice
	disDesc := req.DisDesc
	day := req.Day
	daily := req.Daily
	isRmb := req.IsRmb
	isCion := req.IsCion
	sort := req.Sort
	status := req.Status
	card := models.McVipCard{
		CardName: cardName,
		Price:    price,
		DisRate:  disRate,
		DisPrice: disPrice,
		DisDesc:  disDesc,
		Day:      day,
		Daily:    daily,
		IsRmb:    isRmb,
		IsCion:   isCion,
		Sort:     sort,
		Status:   status,
		Addtime:  utils.GetUnix(),
	}
	if err = global.DB.Create(&card).Error; err != nil {
		return 0, err
	}

	return card.Id, nil
}

func UpdateVipCard(req *models.UpdateCardReq) (res bool, err error) {
	id := req.CardId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	cardName := strings.TrimSpace(req.CardName)
	if cardName == "" {
		err = fmt.Errorf("%v", "会员卡名称不能为空")
		return
	}
	if !CheckCardNameUnique(cardName, id) {
		err = fmt.Errorf("%v", "会员卡名称已经存在")
		return
	}
	price := req.Price
	disRate := req.DisRate
	disPrice := req.DisPrice
	disDesc := req.DisDesc
	status := req.Status
	day := req.Day
	daily := req.Daily
	isRmb := req.IsRmb
	isCion := req.IsCion
	sort := req.Sort
	var mapData = make(map[string]interface{})
	if cardName != "" {
		mapData["card_name"] = cardName
	}
	if price > 0 {
		mapData["price"] = price
	}
	if disRate > 0 {
		mapData["dis_rate"] = disRate
	}
	if disPrice > 0 {
		mapData["dis_price"] = disPrice
	}
	if day > 0 {
		mapData["day"] = day
	}
	if daily > 0 {
		mapData["daily"] = daily
	}
	mapData["dis_desc"] = disDesc
	mapData["is_rmb"] = isRmb
	mapData["is_cion"] = isCion
	mapData["sort"] = sort
	mapData["status"] = status
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McVipCard{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

func DeleteCard(req *models.DeleteCardReq) (res bool, err error) {
	id := req.CardId
	if id <= 0 {
		err = fmt.Errorf("%v", "删除失败，参数错误")
		return

	}
	err = global.DB.Where("id = ?", id).Delete(&models.McVipCard{}).Error
	if err != nil {
		return
	}
	return true, nil
}
