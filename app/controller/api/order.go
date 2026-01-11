package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/order_service"
	"go-novel/global"
	"go-novel/utils"
)

type Order struct{}

func (order *Order) CreateOrder(c *gin.Context) {
	var req models.CreateOrderReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.ClientIp = utils.RemoteIp(c)
	var res *models.UnifiedOrderDataRes
	var err error
	res, err = order_service.CreateOrder(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, res, "生成订单成功")
	return
}

func (order *Order) NotifyUrlTest(c *gin.Context) {
	var req models.OrderNotifyTestReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	var err error
	err = order_service.OrderNotifyTestHandle(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}

	utils.Success(c, "", "订单处理成功")
	return
}

func (order *Order) NotifyUrl(c *gin.Context) {
	var req models.OrderNotifyReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	global.Paylog.Infof("NotifyUrl req=%+v", req)
	err := order_service.OrderNotifyHandle(&req)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	utils.Success(c, "", "订单处理成功")
	return
}

func (order *Order) ReturnUrl(c *gin.Context) {
	// 获取请求参数
	p := c.Request.URL.Query()
	global.Paylog.Infof("ReturnUrl p=%+v", p)
	err := order_service.OrderCallbackHandle(&p)
	if err != nil {
		utils.Fail(c, err, "")
		return
	}
	utils.Success(c, "", "跳转地址")
	return
}

func (order *Order) QueryOrder(c *gin.Context) {
	var req models.QueryOrderReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	var res *models.QueryOrderDataRes
	var err error
	res, err = order_service.QueryOrderHandle(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, res, "查询订单成功")
	return
}
