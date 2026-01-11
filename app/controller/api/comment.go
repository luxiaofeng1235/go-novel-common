package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/models"
	"go-novel/app/service/api/comment_service"
	"go-novel/utils"
)

type Comment struct{}

func (comment *Comment) List(c *gin.Context) {
	var req models.CommentListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	res, err := comment_service.GetCommentRes(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, res, "获取列表成功")
}

func (comment *Comment) Add(c *gin.Context) {
	var req models.CommentAddReq
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
	commentRes, err := comment_service.CommentAdd(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	res := gin.H{
		"comment": commentRes,
	}

	utils.SuccessEncrypt(c, res, "评论成功")
}

func (comment *Comment) Reply(c *gin.Context) {
	var req models.CommentReplyListReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		utils.FailEncrypt(c, err, "参数绑定失败")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	info, total, err := comment_service.GetCommentReplyList(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	res := gin.H{
		"info":  info,
		"total": total,
	}

	utils.SuccessEncrypt(c, res, "ok")
}

func (comment *Comment) Del(c *gin.Context) {
	var req models.CommentDelReq
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
	err = comment_service.CommentDel(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "删除完成")
}

func (comment *Comment) Praise(c *gin.Context) {
	var req models.PraiseUserReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if !ok {
		utils.FailEncrypt(c, nil, "获取登陆用户信息失败")
		return
	}
	req.UserId = userId.(int64)
	err := comment_service.PraiseUser(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	var msg string
	if req.PraiseType > 0 {
		msg = "点赞成功"
	} else {
		msg = "取消点赞成功"
	}
	utils.SuccessEncrypt(c, "", msg)
}

func (comment *Comment) StarGroup(c *gin.Context) {
	var req models.StarGroupReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}
	results, err := comment_service.StarGroup(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, results, "ok")
}

func (comment *Comment) Report(c *gin.Context) {
	var req models.CommentReportReq
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		//参数绑定失败
		utils.FailEncrypt(c, err, "")
		return
	}
	userId, ok := c.Get("user_id")
	if ok {
		req.UserId = userId.(int64)
	}
	err := comment_service.CommentReport(&req)
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}
	utils.SuccessEncrypt(c, "", "ok")
}
