package order_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
)

func GetOrderById(id int64) (order *models.McOrder, err error) {
	err = global.DB.Model(models.McOrder{}).Where("id", id).First(&order).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return order, err
}

func GetOrderByOrderNo(orderNo string) (order *models.McOrder, err error) {
	err = global.DB.Model(models.McOrder{}).Where("order_no = ?", orderNo).First(&order).Error
	return
}

func UpdateOrderPaySuccess(req *models.UpdateOrderSuccessReq) (res bool, err error) {
	if req.OrderId <= 0 {
		err = fmt.Errorf("%s", "orderID不正确")
		return
	}
	id := req.OrderId
	var mapData = make(map[string]interface{})
	mapData["pay_success_time"] = utils.GetUnix()
	mapData["status"] = req.Status
	if req.TradeNo != "" {
		mapData["trade_no"] = req.TradeNo
	}
	if req.TradeAmount > 0 {
		mapData["trade_amount"] = req.TradeAmount
	}
	mapData["uptime"] = utils.GetUnix()
	if err = global.DB.Model(models.McOrder{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}
