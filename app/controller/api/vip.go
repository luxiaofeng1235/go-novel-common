package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/vip_service"
	"go-novel/global"
	"go-novel/utils"
)

type Vip struct{}

func (vip *Vip) GetVipMessage(c *gin.Context) {
	var req models.VipMessageReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	value, err := vip_service.GetVipMessage(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, value, "ok")
}

func (vip *Vip) GetVipCardRmb(c *gin.Context) {
	cards, err := vip_service.GetVipCardRmb()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, cards, "ok")
}

func (vip *Vip) GetVipCardCion(c *gin.Context) {
	cards, err := vip_service.GetVipCardCion()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, cards, "ok")
}

func (vip *Vip) VipBookStore(c *gin.Context) {
	var req models.VipBookStoreReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	req.Ip = utils.RemoteIp(c)
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	isVip, vipCount, choices, hots, news, err := vip_service.GetVipBookStore(&req)
	if err != nil {
		global.Errlog.Infof("记录当前请求游客返回数据： errors=%+v", err)
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"is_vip":    isVip,
		"vip_count": vipCount,
		"choices":   choices,
		"hots":      hots,
		"news":      news,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (vip *Vip) VipBooks(c *gin.Context) {
	var req models.VipChoicesReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	//获取对应的header的mark信息
	mark := utils.GetDeviceQdhInfo(c)
	req.Mark = mark
	req.Ip = utils.RemoteIp(c)
	list, err := vip_service.GetVipBooks(req.IsChoice, req.IsHot, req.IsNew, req.Size, req.DeviceType, req.PackageName, req.Ip, req.Mark)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list": list,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (vip *Vip) ExchangeVip(c *gin.Context) {
	var req models.ExchangeVipReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := vip_service.ExchangeVip(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}
