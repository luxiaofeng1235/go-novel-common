package admin

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/admin/withdraw_service"
	"go-novel/utils"
	"strconv"
)

type Withdraw struct{}

func (withdraw *Withdraw) WithdrawAccountList(c *gin.Context) {
	var req models.WithdrawAccountListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := withdraw_service.WithdrawAccountListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (withdraw *Withdraw) LimitList(c *gin.Context) {
	var req models.WithdrawLimitListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := withdraw_service.LimitListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (withdraw *Withdraw) CreateLimit(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.CreateLimitReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		InsertId, err := withdraw_service.CreateLimit(&req)
		if err != nil {
			utils.Fail(c, err, "创建失败")
			return
		}

		utils.Success(c, InsertId, "ok")
		return
	}

	utils.Success(c, "", "ok")
}

func (withdraw *Withdraw) UpdateLimit(c *gin.Context) {
	if c.Request.Method == "POST" {
		var req models.UpdateLimitReq
		// 参数绑定
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Fail(c, err, "参数绑定失败")
			return
		}
		isUpdate, err := withdraw_service.UpdateLimit(&req)
		if err != nil || isUpdate == false {
			utils.Fail(c, err, "修改信息失败")
			return
		}

		utils.Success(c, "", "ok")
		return
	}

	limitId, _ := strconv.Atoi(c.Query("id"))
	limitInfo, err := withdraw_service.GetLimitById(int64(limitId))
	if err != nil {
		utils.Fail(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"limitInfo": limitInfo,
	}
	utils.Success(c, res, "ok")
}

func (withdraw *Withdraw) DelLimit(c *gin.Context) {
	var req models.DeleteLimitReq
	// 参数绑定
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	isDelete, err := withdraw_service.DeleteLimit(&req)
	if err != nil || isDelete == false {
		utils.Fail(c, err, "删除信息失败")
		return
	}

	utils.Success(c, "", "删除信息成功")
	return
}

func (withdraw *Withdraw) WithdrawList(c *gin.Context) {
	var req models.WithdrawListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}

	list, total, err := withdraw_service.WithdrawListSearch(&req)
	if err != nil {
		utils.Fail(c, err, "获取列表失败")
		return
	}

	res := gin.H{
		"total":       total,
		"currentPage": req.PageNum,
		"list":        list,
	}
	utils.Success(c, res, "ok")
}

func (withdraw *Withdraw) WithdrawCheck(c *gin.Context) {
	var req models.WithdrawCheckReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, err, "参数绑定失败")
		return
	}
	if err := withdraw_service.WithdrawCheck(&req); err != nil {
		utils.Fail(c, err, "审核失败")
		return
	}

	utils.Success(c, "", "")
}
