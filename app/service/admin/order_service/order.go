package order_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/order_service"
	"go-novel/app/service/api/user_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

func GetOrderById(id int64) (order *models.McOrder, err error) {
	err = global.DB.Model(models.McOrder{}).Where("id", id).First(&order).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return order, err
}

func PaySuccessVipCard(user *models.McUser, order *models.McOrder) (err error) {
	cardId := order.CardId
	cardName := order.CardName
	//金币扣款 增加vip 余额改变记录 兑换记录
	day := order.CardDay
	if day > 0 {
		us := make(map[string]interface{})
		us["vip"] = 1
		timeStamp := utils.GetUnix()
		if user.Viptime > timeStamp {
			us["viptime"] = user.Viptime + 86400*day
		} else {
			us["viptime"] = utils.GetUnix() + 86400*day
		}
		err = global.DB.Model(models.McUser{}).Where("id = ?", user.Id).Updates(us).Error
		if err != nil {
			return
		}
	}

	card := models.McVipMessage{
		Uid:      user.Id,
		Nickname: user.Nickname,
		Pic:      user.Pic,
		CardId:   cardId,
		CardName: cardName,
		PayType:  2,
		Addtime:  utils.GetUnix(),
	}
	err = global.DB.Model(models.McVipMessage{}).Create(&card).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func OrderListSearch(req *models.OrderListReq) (list []*models.McOrder, total int64, err error) {
	db := global.DB.Model(&models.McOrder{}).Order("id desc")

	orderNo := strings.TrimSpace(req.OrderNo)
	if orderNo != "" {
		db = db.Where("order_no = ? or  trade_no = ?", orderNo, orderNo)
	}

	userId := strings.TrimSpace(req.UserId)
	if userId != "" {
		db = db.Where("user_id = ?", userId)
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	payType := strings.TrimSpace(req.PayType)
	if payType != "" {
		db = db.Where("pay_type = ?", payType)
	}

	if req.BeginTime != "" {
		db = db.Where("addtime >=?", utils.DateToUnix(req.BeginTime))
	}

	if req.EndTime != "" {
		db = db.Where("addtime <=?", utils.DateToUnix(req.EndTime))
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

func UpdateOrder(req *models.UpdateOrderReq) (res bool, err error) {
	id := req.OrderId
	if id <= 0 {
		err = fmt.Errorf("%v", "id不正确")
		return
	}
	order, err := order_service.GetOrderById(id)
	if err != nil {
		return
	}
	if order.Status != 0 {
		err = fmt.Errorf("%v", "只可以修改代付款的订单")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["status"] = req.Status
	mapData["trade_amount"] = req.TradeAmount
	mapData["pay_success_time"] = utils.GetUnix()
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McOrder{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}

	user, err := user_service.GetUserById(order.UserId)
	if err != nil {
		err = fmt.Errorf("获取用户信息失败 %v", err.Error())
		return
	}
	err = PaySuccessVipCard(user, order)
	return true, nil
}
