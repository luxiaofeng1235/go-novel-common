package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/withdraw_service"
	"go-novel/utils"
)

type Withdraw struct{}

func (withdraw *Withdraw) GetWithdrawLimit(c *gin.Context) {
	limits, err := withdraw_service.GetWithdrawLimit()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, limits, "ok")
}

func (withdraw *Withdraw) AccountDetail(c *gin.Context) {
	var req models.WithdrawAccountDetailReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	alipay, wxpay, err := withdraw_service.GetAccountDetailById(&req)
	if err != nil {
		utils.FailEncrypt(c, nil, "获取数据失败")
		return
	}
	res := gin.H{
		"alipay": alipay,
		"wxpay":  wxpay,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (withdraw *Withdraw) AccountSave(c *gin.Context) {
	var req models.WithdrawAccountSaveReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := withdraw_service.AccountSave(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}

func (withdraw *Withdraw) AccountDel(c *gin.Context) {
	var req models.WithdrawAccountDelReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	var err error
	err = withdraw_service.AccountDel(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, "", "删除完成")
}

func (withdraw *Withdraw) Apply(c *gin.Context) {
	var req models.WithdrawApplyReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := withdraw_service.Apply(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}

func (withdraw *Withdraw) WithdrawList(c *gin.Context) {
	var req models.WithdrawApplyListReq

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, err := withdraw_service.WithdrawApplyList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list":  list,
		"total": total,
	}

	utils.SuccessEncrypt(c, res, "ok")
}
