package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/notice_service"
	"go-novel/utils"
)

type Message struct{}

func (message *Message) GetLastNotice(c *gin.Context) {
	notice, err := notice_service.GetLastNotice()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, notice, "获取最新公告成功")
}

func (message *Message) MessageList(c *gin.Context) {
	var req models.MessageListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	u, _ := c.Get("user_id")
	userId := u.(int64)
	list, total, page, size, err := notice_service.MessageListSearch(&req, userId)
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

func (message *Message) ReplyList(c *gin.Context) {
	var req models.ReplyListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	list, total, page, size, err := notice_service.ReplyList(&req)
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

func (message *Message) PraiseList(c *gin.Context) {
	var req models.PraiseListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}

	list, total, page, size, err := notice_service.PraiseList(&req)
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

func (message *Message) UpdateIsRead(c *gin.Context) {
	var req models.UpdateIsReadReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := notice_service.UpdateIsRead(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}
