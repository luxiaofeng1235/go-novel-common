package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/feedback_service"
	"go-novel/utils"
)

type Feedback struct{}

func (fd *Feedback) HelpList(c *gin.Context) {
	var req models.HelpListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}

	list, total, page, size, err := feedback_service.HelpListSearch(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}

	utils.SuccessEncrypt(c, res, "获取列表成功")
}

func (fd *Feedback) HelpDetail(c *gin.Context) {
	var req models.HelpDetailReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}

	helplInfo, err := feedback_service.GetFeedbackHelpById(req.HelpId)
	if err != nil {
		utils.FailEncrypt(c, nil, "获取数据失败")
		return
	}

	res := gin.H{
		"helplInfo": helplInfo,
	}
	utils.SuccessEncrypt(c, res, "ok")
}

func (fd *Feedback) List(c *gin.Context) {
	var req models.FeedBackListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, page, size, err := feedback_service.List(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}

	utils.SuccessEncrypt(c, res, "获取列表成功")
}

func (fd *Feedback) Add(c *gin.Context) {
	var req models.FeedBackAddReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	req.Ip = utils.RemoteIp(c)
	err := feedback_service.FeedBackAdd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "反馈成功")
}
