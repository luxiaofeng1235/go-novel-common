package order_service

import (
	"encoding/json"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/vip_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"net/url"
)

func CreateOrder(req *models.CreateOrderReq) (unifiedData *models.UnifiedOrderDataRes, err error) {
	cardId := req.Vid
	payType := req.PayType
	userId := req.UserId
	clientIp := req.ClientIp
	if cardId <= 0 {
		err = fmt.Errorf("%v", "会员卡ID不能为空")
		return
	}
	if payType <= 0 {
		err = fmt.Errorf("%v", "支付类型不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	vip, err := vip_service.GetVipById(cardId)
	if err != nil {
		return
	}
	if vip.Id <= 0 {
		err = fmt.Errorf("%v", "会员卡异常,请稍后尝试")
		return
	}
	subject := vip.CardName
	amount := vip.DisPrice
	cardName := vip.CardName
	cardDay := vip.Day
	orderNo := utils.BuildOrderNo("V")
	order := models.McOrder{
		OrderNo:     orderNo,
		Subject:     subject,
		PayType:     payType,
		UserId:      userId,
		CardId:      cardId,
		CardName:    cardName,
		CardDay:     cardDay,
		TotalAmount: amount,
		ClientIp:    clientIp,
		Addtime:     utils.GetUnix(),
		Uptime:      utils.GetUnix(),
	}
	if err = global.DB.Create(&order).Error; err != nil {
		err = fmt.Errorf("创建订单失败，稍后再试 err=%v", err.Error())
		return
	}
	body := cardName
	payAmount := int64(float64(amount * 100))
	var wayCode int = utils.ChannelCode
	if payType == 1 {
		wayCode = utils.AliChannelCode
	} else if payType == 2 {
		wayCode = utils.WxChannelCode
	}
	extParam := ""
	reqTime := fmt.Sprintf("%v", utils.GetUnix())
	params := make(map[string]interface{})
	params["mchId"] = utils.Appid
	params["wayCode"] = wayCode
	params["subject"] = subject
	params["body"] = body
	params["outTradeNo"] = orderNo
	params["amount"] = payAmount
	params["extParam"] = extParam
	params["clientIp"] = clientIp
	params["notifyUrl"] = utils.NotifyUrl
	params["returnUrl"] = utils.ReturnUrl
	params["reqTime"] = reqTime
	sign := utils.PaySign(params, utils.AppSecret)
	params["sign"] = sign
	//log.Printf("%v", params)

	var result string
	result, err = utils.PayPostRequest(utils.UnifiedOrderUrl, params)

	var unRes models.UnifiedOrderRes
	err = json.Unmarshal([]byte(result), &unRes)
	if err != nil {
		return
	}
	if unRes.Code != 0 {
		err = fmt.Errorf("%v", unRes.Message)
	}
	unifiedData = unRes.Data
	return
}

func OrderNotifyTestHandle(req *models.OrderNotifyTestReq) (err error) {
	status := req.Status
	orderNo := req.OrderNo
	tradeAmount := req.TradeAmount
	tradeNo := req.TradeNo
	if status != 1 {
		return
	}
	if orderNo == "" {
		global.Paylog.Errorf("%v", "回调订单号为空")
		return
	}
	order, err := GetOrderByOrderNo(orderNo)
	if err != nil {
		global.Paylog.Errorf("订单异常 %v", err.Error())
		err = fmt.Errorf("%v", "订单异常")
		return
	}
	if order.Id <= 0 {
		global.Paylog.Errorf("orderNo=%v不存在", orderNo)
		err = fmt.Errorf("%v", "订单不存在")
		return
	}
	if order.Status != 0 {
		global.Paylog.Errorf("orderNo=%v 订单状态异常", order.OrderNo)
		err = fmt.Errorf("%v", "订单状态异常")
		return
	}
	totalAmount := order.TotalAmount
	if totalAmount != tradeAmount {
		global.Paylog.Errorf("orderNo=%v 订单金额异常 TotalAmount=%v tradeAmount=%v", order.OrderNo, totalAmount, tradeAmount)
		err = fmt.Errorf("%v", "订单金额异常")
		return
	}
	orderReq := models.UpdateOrderSuccessReq{
		OrderId:        order.Id,
		TradeNo:        tradeNo,
		TradeAmount:    totalAmount,
		Status:         1,
		PaySuccessTime: utils.GetUnix(),
	}
	var isUpdate bool
	isUpdate, err = UpdateOrderPaySuccess(&orderReq)
	if isUpdate == false || err != nil {
		err = fmt.Errorf("订单状态更改失败 %v", err.Error())
		return
	}
	return
}

func OrderNotifyHandle(req *models.OrderNotifyReq) (err error) {
	mchId := req.MchId
	tradeNo := req.TradeNo
	orderNo := req.OutTradeNo
	originTradeNo := req.OriginTradeNo
	amount := req.Amount
	subject := req.Subject
	body := req.Body
	extParam := req.ExtParam
	state := req.State
	notifyTime := req.NotifyTime
	sign := req.Sign

	params := make(map[string]interface{})
	params["mchId"] = mchId
	params["tradeNo"] = tradeNo
	params["outTradeNo"] = orderNo
	params["originTradeNo"] = originTradeNo
	params["amount"] = amount
	params["subject"] = subject
	params["body"] = body
	params["extParam"] = extParam
	params["state"] = state
	params["notifyTime"] = notifyTime
	dataSign := utils.PaySign(params, utils.AppSecret)
	if sign != dataSign {
		err = fmt.Errorf("%v", "验签失败")
		global.Paylog.Errorf("验签失败 params=%v sign=%v dataSign=%v", params, sign, dataSign)
		return
	}

	order, err := GetOrderByOrderNo(orderNo)
	if err != nil {
		global.Paylog.Errorf("订单异常 %v", err.Error())
		err = fmt.Errorf("%v", "订单异常")
		return
	}
	if order.Id <= 0 {
		global.Paylog.Errorf("orderNo=%v不存在", orderNo)
		err = fmt.Errorf("%v", "订单不存在")
		return
	}
	if order.Status != 0 {
		global.Paylog.Errorf("orderNo=%v 订单状态异常", order.OrderNo)
		err = fmt.Errorf("%v", "订单状态异常")
		return
	}
	totalAmount := order.TotalAmount
	payAmount := float64(amount) / 100
	if totalAmount != payAmount {
		global.Paylog.Errorf("orderNo=%v 订单金额异常 TotalAmount=%v payAmount=%v", order.OrderNo, totalAmount, payAmount)
		err = fmt.Errorf("%v", "订单金额异常")
		return
	}
	orderReq := models.UpdateOrderSuccessReq{
		OrderId:        order.Id,
		TradeNo:        tradeNo,
		TradeAmount:    totalAmount,
		Status:         1,
		PaySuccessTime: utils.GetUnix(),
	}
	var isUpdate bool
	isUpdate, err = UpdateOrderPaySuccess(&orderReq)
	if isUpdate == false || err != nil {
		err = fmt.Errorf("订单状态更改失败 %v", err.Error())
		return
	}

	userId := order.UserId

	user, err := getUserById(userId)
	if err != nil {
		err = fmt.Errorf("获取用户信息失败 %v", err.Error())
		return
	}
	err = PaySuccessVipCard(user, order)
	return
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

func OrderCallbackHandle(val *url.Values) (err error) {
	log.Printf("跳转参数:%v", val)
	orderNo := val.Get("outTradeNo")

	order, err := GetOrderByOrderNo(orderNo)
	if err != nil {
		err = fmt.Errorf(" 获取订单数据失败:%v ", err)
		//global.NhPaylog.Infof("订单金额转换失败 :%v parames=%v", err.Error(), p)
		return
	}
	if order != nil && order.Id <= 0 {
		err = fmt.Errorf("%v 获取订单数据失败:", err)
		global.Paylog.Infof("获取订单数据失败 :%v parames=%v", err.Error(), val)
		return
	}
	return
}

func QueryOrderHandle(req *models.QueryOrderReq) (queryData *models.QueryOrderDataRes, err error) {
	orderNo := req.OrderNo
	if orderNo == "" {
		err = fmt.Errorf("%v", "订单号不能为空")
		return
	}
	order, err := GetOrderByOrderNo(orderNo)
	if err != nil {
		global.Paylog.Errorf("订单异常 %v", err.Error())
		err = fmt.Errorf("订单异常 %v", err.Error())
		return
	}
	if order.Id <= 0 {
		global.Paylog.Errorf("orderNo=%v不存在", orderNo)
		err = fmt.Errorf("%v", "订单不存在")
		return
	}

	params := make(map[string]interface{})
	params["mchId"] = utils.Appid
	params["outTradeNo"] = orderNo
	params["reqTime"] = utils.GetUnix()
	sign := utils.PaySign(params, utils.AppSecret)
	params["sign"] = sign
	//log.Printf("%v", params)

	var result string
	result, err = utils.PayPostRequest(utils.QueryOrderUrl, params)
	log.Println(result)
	var unRes models.QueryOrderRes
	err = json.Unmarshal([]byte(result), &unRes)
	if err != nil {
		return
	}
	log.Println(unRes)
	if unRes.Code != 0 {
		err = fmt.Errorf("%v", unRes.Message)
	}
	queryData = unRes.Data
	return
}
